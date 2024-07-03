package service

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/rarimo/forms-svc/internal/config"
	"github.com/rarimo/forms-svc/internal/data/pg"
	"github.com/rarimo/forms-svc/internal/service/handlers"
	"gitlab.com/distributed_lab/ape"
)

func Run(ctx context.Context, cfg config.Config) {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(cfg.Log()),
		ape.LoganMiddleware(cfg.Log()),
		ape.CtxMiddleware(
			handlers.CtxLog(cfg.Log()),
			handlers.CtxFormsQ(pg.NewForms(cfg.DB().Clone())),
			handlers.CtxForms(cfg.Forms()),
		),
	)
	r.Route("/integrations/forms-svc/v1", func(r chi.Router) {
		r.Route("/form", func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(cfg.Auth(), cfg.Log()))
			r.Post("/submit", handlers.SubmitForm)
			r.Get("/{id}", handlers.GetForm)
		})
	})

	cfg.Log().Info("Service started")
	ape.Serve(ctx, r, cfg, ape.ServeOpts{})
}
