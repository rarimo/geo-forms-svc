package config

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rarimo/geo-forms-svc/internal/data"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

const formsTable = "consulate"

type Forms struct {
	Cooldown          time.Duration
	Period            time.Duration
	MinAbnormalPeriod time.Duration
	MaxAbnormalPeriod time.Duration
	ResendFormsCount  uint64
	db                *sql.DB
}

func (c *config) Forms() *Forms {
	return c.forms.Do(func() interface{} {
		var cfg struct {
			Cooldown          time.Duration `fig:"cooldown,required"`
			Period            time.Duration `fig:"period,required"`
			MinAbnormalPeriod time.Duration `fig:"min_abnormal_period,required"`
			MaxAbnormalPeriod time.Duration `fig:"max_abnormal_period,required"`
			ResendFormsCount  uint64        `fig:"resend_forms_count,required"`
			URL               string        `fig:"url,required"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "forms")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out withdrawal point price: %w", err))
		}

		db, err := sql.Open("mysql", cfg.URL)
		if err != nil {
			panic(fmt.Errorf("failed to connect to mysql: %w", err))
		}

		return &Forms{
			Cooldown:          cfg.Cooldown,
			Period:            cfg.Period,
			MinAbnormalPeriod: cfg.MinAbnormalPeriod,
			MaxAbnormalPeriod: cfg.MaxAbnormalPeriod,
			ResendFormsCount:  cfg.ResendFormsCount,
			db:                db,
		}
	}).(*Forms)
}

// SendForms inserts the forms to another remote database, which URL is provided in config.
func (f *Forms) SendForms(forms ...*data.Form) error {
	if len(forms) == 0 {
		return nil
	}

	stmt := squirrel.Insert(formsTable).Columns(
		"name",
		"surname",
		"id_num",
		"birthday",
		"citizen",
		"visited",
		"purpose",
		"country",
		"city",
		"address",
		"postal",
		"phone",
		"email",
		"image",
	)

	for _, form := range forms {
		if form == nil {
			continue
		}
		stmt = stmt.Values(
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
			form.Image,
		)
	}

	query, args, err := stmt.ToSql()
	if err != nil {
		return fmt.Errorf("failed to construct db query: %w", err)
	}

	if _, err = f.db.Exec(query, args...); err != nil {
		return fmt.Errorf("insert form: %w", err)
	}

	return nil
}