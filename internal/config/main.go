package config

import (
	"github.com/rarimo/geo-auth-svc/pkg/auth"
	"github.com/rarimo/geo-forms-svc/internal/service/workers/spreadsheets"
	"github.com/rarimo/geo-forms-svc/internal/storage"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type Config interface {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	auth.Auther
	storage.Storager
	spreadsheets.Spreadsheeter

	Forms() *Forms
}

type config struct {
	comfig.Logger
	pgdb.Databaser
	comfig.Listenerer
	auth.Auther
	storage.Storager
	spreadsheets.Spreadsheeter

	forms comfig.Once

	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:        getter,
		Databaser:     pgdb.NewDatabaser(getter),
		Listenerer:    comfig.NewListenerer(getter),
		Logger:        comfig.NewLogger(getter, comfig.LoggerOpts{}),
		Auther:        auth.NewAuther(getter),
		Storager:      storage.NewStorager(getter),
		Spreadsheeter: spreadsheets.NewSpreadsheeter(getter),
	}
}
