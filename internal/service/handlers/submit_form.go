package handlers

import (
	"database/sql"
	"net/http"
	"net/url"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/geo-forms-svc/internal/data"
	"github.com/rarimo/geo-forms-svc/internal/service/requests"
	"github.com/rarimo/geo-forms-svc/internal/storage"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func SubmitForm(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewSubmitForm(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	nullifier := strings.ToLower(UserClaims(r)[0].Nullifier)

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

	imageURL, err := url.Parse(req.Data.Attributes.Image)
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to parse image URL %s", req.Data.Attributes.Image)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if err = Storage(r).ValidateImage(imageURL); err != nil {
		if storage.IsBadRequestError(err) {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"image": err,
			})...)
			return
		}

		Log(r).WithError(err).Error("Failed to validate image")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	userData := req.Data.Attributes
	form := &data.Form{
		Nullifier: nullifier,
		Status:    data.AcceptedStatus,
		Name:      userData.Name,
		Surname:   userData.Surname,
		IDNum:     userData.IdNum,
		Birthday:  userData.Birthday,
		Citizen:   userData.Citizen,
		Visited:   userData.Visited,
		Purpose:   userData.Purpose,
		Country:   userData.Country,
		City:      userData.City,
		Address:   userData.Address,
		Postal:    userData.Postal,
		Phone:     userData.Phone,
		Email:     userData.Email,
		Image:     nil,
		ImageURL:  sql.NullString{String: userData.Image, Valid: true},
	}

	formStatus, err := FormsQ(r).Insert(form)
	if err != nil {
		Log(r).WithError(err).Error("failed to insert form")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	formStatus.NextFormAt = formStatus.CreatedAt.Add(Forms(r).Cooldown)

	ape.Render(w, newFormStatusResponse(formStatus))
}
