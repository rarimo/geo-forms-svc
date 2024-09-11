package handlers

import (
	"net/http"

	"github.com/rarimo/geo-forms-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetImage(w http.ResponseWriter, r *http.Request) {
	id, apiKey, err := requests.NewGetImage(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if Storage(r).APIKey != apiKey {
		Log(r).Warnf("Request with apiKey=%s, but want: %s", apiKey, Storage(r).APIKey)
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	signedURL, err := Storage(r).RawSignedGetURL(id)
	if err != nil {
		Log(r).WithError(err).Error("Failed to sign url")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.Header().Set("Location", signedURL)
	w.WriteHeader(http.StatusFound)
}
