package handlers

import (
	"net/http"

	"github.com/rarimo/geo-auth-svc/pkg/auth"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3"
)

func AuthMiddleware(auth *auth.Client, log *logan.Entry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := auth.ValidateJWT(r)
			if err != nil {
				log.WithError(err).Info("Got invalid auth or validation error")
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			if len(claims) == 0 {
				ape.RenderErr(w, problems.Unauthorized())
				return
			}

			ctx := CtxUserClaims(claims)(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
