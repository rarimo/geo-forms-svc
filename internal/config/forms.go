package config

import (
	"fmt"
	"time"

	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/kv"
)

type Forms struct {
	Cooldown time.Duration
}

func (c *config) Forms() *Forms {
	return c.forms.Do(func() interface{} {
		var cfg struct {
			Cooldown time.Duration `fig:"cooldown,required"`
		}

		err := figure.Out(&cfg).
			From(kv.MustGetStringMap(c.getter, "forms")).
			Please()
		if err != nil {
			panic(fmt.Errorf("failed to figure out withdrawal point price: %w", err))
		}

		return &Forms{
			Cooldown: cfg.Cooldown,
		}
	}).(*Forms)
}
