package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rarimo/geo-auth-svc/pkg/auth"
	"github.com/rarimo/geo-forms-svc/internal/data"
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

	lastForm, err := FormsQ(r).FilterByNullifier(nullifier).Last()
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to get last user form for nullifier [%s]", nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if lastForm == nil {
		signedURL, id, err := newCreatedFormWithURL(r, nullifier, req.Data.Attributes.ContentType, req.Data.Attributes.ContentLength)
		if err != nil {
			Log(r).WithError(err).Error("Failed to create form")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, newUploadImageResponse(id, signedURL))
		return
	}

	if lastForm.Status == data.CreatedStatus {
		signedURL, id, err := Storage(r).GeneratePutURL(lastForm.ID, req.Data.Attributes.ContentType, req.Data.Attributes.ContentLength)
		if err != nil {
			Log(r).WithError(err).Error("Failed to generate pre-signed url")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, newUploadImageResponse(id, signedURL))
		return
	}

	next := lastForm.CreatedAt.Add(Forms(r).Cooldown)
	if next.After(time.Now().UTC()) {
		Log(r).Debugf("Form submitted time: %s; next available time: %s", lastForm.CreatedAt, next)
		ape.RenderErr(w, problems.TooManyRequests())
		return
	}

	signedURL, id, err := newCreatedFormWithURL(r, nullifier, req.Data.Attributes.ContentType, req.Data.Attributes.ContentLength)
	if err != nil {
		Log(r).WithError(err).Error("Failed to create form")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newUploadImageResponse(id, signedURL))
}

func newUploadImageResponse(id, signedURL string) resources.UploadImageResponseResponse {
	return resources.UploadImageResponseResponse{
		Data: resources.UploadImageResponse{
			Key: resources.Key{
				ID:   id,
				Type: resources.UPLOAD_IMAGE_RESPONSE,
			},
			Attributes: resources.UploadImageResponseAttributes{
				Url: signedURL,
			},
		},
	}
}

func newCreatedFormWithURL(r *http.Request, nullifier, contentType string, contentLength int64) (string, string, error) {
	signedURL, id, err := Storage(r).GeneratePutURL("", contentType, contentLength)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate pre-signed url: %w", err)
	}

	err = FormsQ(r).Insert(data.Form{
		ID:        id,
		Status:    data.CreatedStatus,
		Nullifier: nullifier,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to insert created form: %w", err)
	}

	return signedURL, id, nil
}
