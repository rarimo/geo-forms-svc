package formsender

import (
	"context"
	"fmt"

	"github.com/rarimo/geo-forms-svc/internal/config"
	"github.com/rarimo/geo-forms-svc/internal/data"
	"github.com/rarimo/geo-forms-svc/internal/data/pg"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/running"
)

type dbaser struct {
	db *pgdb.DB
}

func Run(ctx context.Context, cfg config.Config) {
	log := cfg.Log().WithField("who", "form-sender")
	formsCfg := cfg.Forms()
	db := dbaser{cfg.DB().Clone()}

	running.WithBackOff(
		ctx,
		log,
		"resender",
		func(context.Context) error {
			forms, err := db.FormsQ().FilterByStatus(data.AcceptedStatus).Select()
			if err != nil {
				return fmt.Errorf("failed to get unsended forms: %w", err)
			}
			if len(forms) == 0 {
				return nil
			}

			if err = cfg.Forms().SendForms(forms...); err != nil {
				return fmt.Errorf("failed to send forms: %w", err)
			}

			ids := make([]string, len(forms))
			for i, v := range forms {
				ids[i] = v.ID
			}

			if err = db.FormsQ().FilterByID(ids...).Update(data.ProcessedStatus); err != nil {
				return fmt.Errorf("failed to update form status: %w", err)
			}

			return nil
		},
		formsCfg.Period,
		formsCfg.MinAbnormalPeriod,
		formsCfg.MaxAbnormalPeriod,
	)
}

func (d *dbaser) FormsQ() data.FormsQ {
	return pg.NewForms(d.db)
}
