package service

import (
	"context"

	"github.com/go-chi/chi"
	"github.com/rarimo/geo-forms-svc/internal/config"
	"github.com/rarimo/geo-forms-svc/internal/data/pg"
	"github.com/rarimo/geo-forms-svc/internal/service/handlers"
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
			handlers.CtxStorage(cfg.Storage()),
		),
	)
	r.Route("/integrations/geo-forms-svc/v1", func(r chi.Router) {
		r.Get("/image/{id}", nil)
		r.With(handlers.AuthMiddleware(cfg.Auth(), cfg.Log())).Post("/image", handlers.UploadImage)
		r.Route("/status", func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(cfg.Auth(), cfg.Log()))
			r.Get("/{id}", handlers.StatusByID)
			r.Get("/last", handlers.LastStatus)
		})
		r.Route("/form", func(r chi.Router) {
			r.Use(handlers.AuthMiddleware(cfg.Auth(), cfg.Log()))
			r.Post("/submit", handlers.LegacySubmitForm)
			r.Post("/", handlers.SubmitForm)
		})
	})

	r.Route("/integrations/geo-forms-svc/v2", func(r chi.Router) {
		r.With(handlers.AuthMiddleware(cfg.Auth(), cfg.Log())).Post("/image", handlers.UploadImageV2)
	})

	cfg.Log().Info("Service started")
	ape.Serve(ctx, r, cfg, ape.ServeOpts{})
}
