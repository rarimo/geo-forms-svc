package config

import (
	"fmt"
	"time"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type Forms struct {
	Cooldown          time.Duration
	Period            time.Duration
	MinAbnormalPeriod time.Duration
	MaxAbnormalPeriod time.Duration
	ResendFormsCount  uint64
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

		return &Forms{
			Cooldown:          cfg.Cooldown,
			Period:            cfg.Period,
			MinAbnormalPeriod: cfg.MinAbnormalPeriod,
			MaxAbnormalPeriod: cfg.MaxAbnormalPeriod,
			ResendFormsCount:  cfg.ResendFormsCount,
		}
	}).(*Forms)
}
