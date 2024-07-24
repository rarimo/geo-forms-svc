package formsender

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rarimo/geo-forms-svc/internal/config"
	"github.com/rarimo/geo-forms-svc/internal/data"
	"github.com/rarimo/geo-forms-svc/internal/data/pg"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/running"
)

type formsQ struct {
	db *pgdb.DB
}

func Run(ctx context.Context, cfg config.Config) {
	log := cfg.Log().WithField("who", "form-sender")
	db := formsQ{cfg.DB().Clone()}
	storage := cfg.Storage()

	running.WithBackOff(ctx, log, "resender", func(context.Context) error {
		forms, err := db.FormsQ().FilterByStatus(data.AcceptedStatus).Limit(cfg.Forms().ResendFormsCount).Select()
		if err != nil {
			return fmt.Errorf("failed to get unsended forms: %w", err)
		}
		if len(forms) == 0 {
			return nil
		}

		for i, form := range forms {
			if form.Image != nil {
				continue
			}

			imageURL, err := url.Parse(form.ImageURL.String)
			if err != nil {
				return fmt.Errorf("failed to parse image url: %w", err)
			}

			forms[i].Image, err = storage.GetImageBase64(imageURL)
			if err != nil {
				return fmt.Errorf("failed to get image base64: %w", err)
			}
		}

		if err = cfg.Forms().SendForms(forms...); err != nil {
			return fmt.Errorf("failed to send forms: %w", err)
		}

		ids := make([]string, len(forms))
		for i, v := range forms {
			ids[i] = v.ID
		}

		err = db.FormsQ().FilterByID(ids...).Update(map[string]any{
			data.ColStatus: data.ProcessedStatus,
		})
		if err != nil {
			return fmt.Errorf("failed to update form status: %w", err)
		}

		return nil
	},
		cfg.Forms().Period,
		cfg.Forms().MinAbnormalPeriod,
		cfg.Forms().MaxAbnormalPeriod,
	)
}

func (d *formsQ) FormsQ() data.FormsQ {
	return pg.NewForms(d.db)
}
