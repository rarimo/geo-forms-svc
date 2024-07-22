package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/rarimo/geo-auth-svc/pkg/auth"
	"github.com/rarimo/geo-forms-svc/internal/service/requests"
	"github.com/rarimo/geo-forms-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewUploadImage(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	nullifier := strings.ToLower(UserClaims(r)[0].Nullifier)
	if !auth.Authenticates(UserClaims(r), auth.VerifiedGrant(nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	lastForm, err := FormsQ(r).Last(nullifier)
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to get last user form for nullifier [%s]", nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if lastForm != nil {
		next := lastForm.CreatedAt.Add(Forms(r).Cooldown)
		if next.After(time.Now().UTC()) {
			Log(r).Debugf("Form submitted time: %s; next available time: %s", lastForm.CreatedAt, next)
			ape.RenderErr(w, problems.TooManyRequests())
			return
		}
	}

	signedURL, key, err := Storage(r).GeneratePUTURL(req.Data.Attributes.ContentType, req.Data.Attributes.ContentLength)
	if err != nil {
		Log(r).WithError(err).Error("Failed to generate pre-signed url")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, resources.UploadImageResponseResponse{
		Data: resources.UploadImageResponse{
			Key: resources.Key{
				ID:   key,
				Type: resources.UPLOAD_IMAGE_RESPONSE,
			},
			Attributes: resources.UploadImageResponseAttributes{
				Url: signedURL,
			},
		},
	})
}
