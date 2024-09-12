package handlers

import (
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

func UploadImageV2(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewUploadImageV2(r)
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

	selfieImage := req.Data.Attributes.SelfieImage
	passportImage := req.Data.Attributes.PassportImage

	var selfieImageSignedURL, passportImageSignedURL, id string

	if lastForm == nil {
		selfieImageSignedURL, id, err = Storage(r).GeneratePutURL("", selfieImage.ContentType, selfieImage.ContentLength)
		if err != nil {
			Log(r).WithError(err).Error("Failed to generate selfie image pre-signed url")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		if passportImage != nil {
			passportImageSignedURL, _, err = Storage(r).GeneratePutURL(id+"-pass", passportImage.ContentType, passportImage.ContentLength)
			if err != nil {
				Log(r).WithError(err).Error("Failed to generate passport image pre-signed url")
				ape.RenderErr(w, problems.InternalError())
				return
			}
		}

		err = FormsQ(r).Insert(data.Form{
			ID:        id,
			Status:    data.CreatedStatus,
			Nullifier: nullifier,
		})
		if err != nil {
			Log(r).WithError(err).Error("Failed to insert created form")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, newUploadImageResponseV2(id, selfieImageSignedURL, passportImageSignedURL))
		return
	}

	if lastForm.Status == data.CreatedStatus {
		selfieImageSignedURL, id, err = Storage(r).GeneratePutURL(lastForm.ID, selfieImage.ContentType, selfieImage.ContentLength)
		if err != nil {
			Log(r).WithError(err).Error("Failed to generate selfie image pre-signed url")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		if passportImage != nil {
			passportImageSignedURL, _, err = Storage(r).GeneratePutURL(id+"-pass", passportImage.ContentType, passportImage.ContentLength)
			if err != nil {
				Log(r).WithError(err).Error("Failed to generate passport image pre-signed url")
				ape.RenderErr(w, problems.InternalError())
				return
			}
		}

		ape.Render(w, newUploadImageResponseV2(id, selfieImageSignedURL, passportImageSignedURL))
		return
	}

	next := lastForm.CreatedAt.Add(Forms(r).Cooldown)
	if next.After(time.Now().UTC()) {
		Log(r).Debugf("Form submitted time: %s; next available time: %s", lastForm.CreatedAt, next)
		ape.RenderErr(w, problems.TooManyRequests())
		return
	}

	selfieImageSignedURL, id, err = Storage(r).GeneratePutURL("", selfieImage.ContentType, selfieImage.ContentLength)
	if err != nil {
		Log(r).WithError(err).Error("Failed to generate selfie image pre-signed url")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if passportImage != nil {
		passportImageSignedURL, _, err = Storage(r).GeneratePutURL(id+"-pass", passportImage.ContentType, passportImage.ContentLength)
		if err != nil {
			Log(r).WithError(err).Error("Failed to generate passport image pre-signed url")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	err = FormsQ(r).Insert(data.Form{
		ID:        id,
		Status:    data.CreatedStatus,
		Nullifier: nullifier,
	})
	if err != nil {
		Log(r).WithError(err).Error("Failed to insert created form")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newUploadImageResponseV2(id, selfieImageSignedURL, passportImageSignedURL))
}

func newUploadImageResponseV2(id, selfieImageSignedURL, passportImageSignedURL string) resources.UploadImageResponseV2Response {
	var passportImageURL *string
	if passportImageSignedURL != "" {
		passportImageURL = &passportImageSignedURL
	}
	return resources.UploadImageResponseV2Response{
		Data: resources.UploadImageResponseV2{
			Key: resources.Key{
				ID:   id,
				Type: resources.UPLOAD_IMAGE_RESPONSE,
			},
			Attributes: resources.UploadImageResponseV2Attributes{
				SelfieImageUrl:   selfieImageSignedURL,
				PassportImageUrl: passportImageURL,
			},
		},
	}
}
