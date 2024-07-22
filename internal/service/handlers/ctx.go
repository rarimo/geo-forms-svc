package handlers

import (
	"context"
	"net/http"

	"github.com/rarimo/geo-auth-svc/resources"
	"github.com/rarimo/geo-forms-svc/internal/config"
	"github.com/rarimo/geo-forms-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	userClaimsCtxKey
	formsQCtxKey
	formsCtxKey
	storageCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxUserClaims(claim []resources.Claim) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, userClaimsCtxKey, claim)
	}
}

func UserClaims(r *http.Request) []resources.Claim {
	return r.Context().Value(userClaimsCtxKey).([]resources.Claim)
}

func CtxFormsQ(q data.FormsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, formsQCtxKey, q)
	}
}

func FormsQ(r *http.Request) data.FormsQ {
	return r.Context().Value(formsQCtxKey).(data.FormsQ).New()
}

func CtxForms(cfg *config.Forms) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, formsCtxKey, cfg)
	}
}

func Forms(r *http.Request) *config.Forms {
	return r.Context().Value(formsCtxKey).(*config.Forms)
}

func CtxStorage(cfg *config.Storage) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, storageCtxKey, cfg)
	}
}

func Storage(r *http.Request) *config.Storage {
	return r.Context().Value(storageCtxKey).(*config.Storage)
}
