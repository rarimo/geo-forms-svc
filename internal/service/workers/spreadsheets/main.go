package spreadsheets

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rarimo/geo-forms-svc/internal/data"
	"github.com/rarimo/geo-forms-svc/internal/data/pg"
	"github.com/rarimo/geo-forms-svc/internal/storage"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/running"
)

type formsQ struct {
	db *pgdb.DB
}

type extConfig interface {
	comfig.Logger
	pgdb.Databaser
	Spreadsheeter
	storage.Storager
}

func Run(ctx context.Context, cfg extConfig) {
	log := cfg.Log().WithField("who", "spreadsheeter")
	db := formsQ{cfg.DB().Clone()}
	s3 := cfg.Storage()
	spreadsheets := cfg.Spreadsheets()

	running.WithBackOff(ctx, log, "sheet-former", func(context.Context) error {
		forms, err := db.FormsQ().FilterByStatus(data.AcceptedStatus).Select()
		if err != nil {
			return fmt.Errorf("failed to get unsended forms: %w", err)
		}
		if len(forms) == 0 {
			return nil
		}

		tableData := make([][]any, 0, len(forms))
		for _, form := range forms {
			data := make([]any, 0, len(headers))

			link, err := url.Parse(form.Image)
			if err != nil {
				return fmt.Errorf("failed to parse image url %s: %w", form.Image, err)
			}

			_, key, err := s3.BucketAndKey(link)
			if err != nil {
				return fmt.Errorf("failed to generate get bucket and key: %w", err)
			}

			selfieImageURL, err := url.Parse(spreadsheets.Self)
			if err != nil {
				return fmt.Errorf("failed to parse self url %s: %w", spreadsheets.Self, err)
			}

			selfieImageURL = selfieImageURL.JoinPath(key)
			q := selfieImageURL.Query()
			q.Add("api", s3.APIKey)
			selfieImageURL.RawQuery = q.Encode()

			var passportImageSignedURL string
			if form.PassportImage.Valid {
				link, err := url.Parse(form.PassportImage.String)
				if err != nil {
					return fmt.Errorf("failed to parse image url %s: %w", form.Image, err)
				}

				_, key, err := s3.BucketAndKey(link)
				if err != nil {
					return fmt.Errorf("failed to generate get bucket and key: %w", err)
				}

				passportImageURL, err := url.Parse(spreadsheets.Self)
				if err != nil {
					return fmt.Errorf("failed to parse self url %s: %w", spreadsheets.Self, err)
				}

				passportImageURL = passportImageURL.JoinPath(key)
				q := passportImageURL.Query()
				q.Add("api", s3.APIKey)
				passportImageURL.RawQuery = q.Encode()
				passportImageSignedURL = passportImageURL.String()
			}

			data = append(data,
				form.Name,
				form.Surname,
				form.IDNum,
				form.Birthday,
				form.Citizen,
				form.Visited,
				form.Purpose,
				form.Country,
				form.City,
				form.Address,
				form.Postal,
				form.Phone,
				form.Email,
				form.UpdatedAt.Format("01/02/2006 15:04"),
				selfieImageURL.String(),
				passportImageSignedURL,
			)

			tableData = append(tableData, data)

			if form.UpdatedAt.After(spreadsheets.lastSubmited) {
				spreadsheets.lastSubmited = form.UpdatedAt
			}

		}

		err = spreadsheets.CreateTable()
		if err != nil {
			return fmt.Errorf("failed to create spreadsheet: %w", err)
		}

		err = spreadsheets.FillTable(tableData)
		if err != nil {
			return fmt.Errorf("failed to fill spreadsheet: %w", err)
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
		spreadsheets.period,
		spreadsheets.minAbnormalPeriod,
		spreadsheets.maxAbnormalPeriod,
	)
}

func (d *formsQ) FormsQ() data.FormsQ {
	return pg.NewForms(d.db)
}
