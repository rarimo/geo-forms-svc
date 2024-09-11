package handlers

import (
	"database/sql"
	"net/http"
	"net/url"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/geo-auth-svc/pkg/auth"
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
	if !auth.Authenticates(UserClaims(r), auth.VerifiedGrant(nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	lastForm, err := FormsQ(r).FilterByNullifier(nullifier).Last()
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to get last user form by nullifier [%s]", nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if lastForm == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if lastForm.Status != data.CreatedStatus {
		Log(r).Debugf("User last form must have %s status, got %s", data.CreatedStatus, lastForm.Status)
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	selfieImageURL, err := url.Parse(req.Data.Attributes.Image)
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to parse selfie image URL %s", req.Data.Attributes.Image)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if err = Storage(r).ValidateImage(selfieImageURL, lastForm.ID); err != nil {
		if storage.IsBadRequestError(err) {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"image": err,
			})...)
			return
		}

		Log(r).WithError(err).Error("Failed to validate selfie image")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var passportImageURL *url.URL
	if req.Data.Attributes.PassportImage != nil {
		passportImageURL, err = url.Parse(*req.Data.Attributes.PassportImage)
		if err != nil {
			Log(r).WithError(err).Errorf("Failed to parse passport image URL %s", *req.Data.Attributes.PassportImage)
			ape.RenderErr(w, problems.InternalError())
			return
		}

		if err = Storage(r).ValidateImage(passportImageURL, lastForm.ID+"-pass"); err != nil {
			if storage.IsBadRequestError(err) {
				ape.RenderErr(w, problems.BadRequest(validation.Errors{
					"passport_image": err,
				})...)
				return
			}

			Log(r).WithError(err).Error("Failed to validate passport image")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	var passportImage sql.NullString
	if req.Data.Attributes.PassportImage != nil {
		passportImage.String = *req.Data.Attributes.PassportImage
		passportImage.Valid = true
	}

	userData := req.Data.Attributes
	err = FormsQ(r).FilterByID(lastForm.ID).Update(map[string]interface{}{
		data.ColStatus:        data.AcceptedStatus,
		data.ColName:          userData.Name,
		data.ColSurname:       userData.Surname,
		data.ColIDNum:         userData.IdNum,
		data.ColBirthday:      userData.Birthday,
		data.ColCitizen:       userData.Citizen,
		data.ColVisited:       userData.Visited,
		data.ColPurpose:       userData.Purpose,
		data.ColCountry:       userData.Country,
		data.ColCity:          userData.City,
		data.ColAddress:       userData.Address,
		data.ColPostal:        userData.Postal,
		data.ColPhone:         userData.Phone,
		data.ColEmail:         userData.Email,
		data.ColImage:         userData.Image,
		data.ColPassportImage: passportImage,
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to update form")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	lastForm, err = FormsQ(r).FilterByNullifier(nullifier).Last()
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to get last user form for nullifier [%s]", nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	nextFormAt := lastForm.CreatedAt.Add(Forms(r).Cooldown)

	ape.Render(w, newFormStatusResponse(*lastForm, nextFormAt))
}
