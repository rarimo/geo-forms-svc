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

	lastForm, err := FormsQ(r).Last(nullifier)
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to get last user form for nullifier [%s]", nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if lastForm == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if lastForm.Status != data.CreatedStatus {
		Log(r).Debugf("User last form don't have created status")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	imageURL, err := url.Parse(req.Data.Attributes.Image)
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to parse image URL %s", req.Data.Attributes.Image)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if err = Storage(r).ValidateImage(imageURL, lastForm.ID); err != nil {
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
	err = FormsQ(r).FilterByID(lastForm.ID).Update(map[string]interface{}{
		data.ColStatus:   data.AcceptedStatus,
		data.ColName:     userData.Name,
		data.ColSurname:  userData.Surname,
		data.ColIDNum:    userData.IdNum,
		data.ColBirthday: userData.Birthday,
		data.ColCitizen:  userData.Citizen,
		data.ColVisited:  userData.Visited,
		data.ColPurpose:  userData.Purpose,
		data.ColCountry:  userData.Country,
		data.ColCity:     userData.City,
		data.ColAddress:  userData.Address,
		data.ColPostal:   userData.Postal,
		data.ColPhone:    userData.Phone,
		data.ColEmail:    userData.Email,
		data.ColImageURL: sql.NullString{String: userData.Image, Valid: true},
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to insert form")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	lastForm, err = FormsQ(r).Last(nullifier)
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to get last user form for nullifier [%s]", nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	lastForm.NextFormAt = lastForm.CreatedAt.Add(Forms(r).Cooldown)

	ape.Render(w, newFormStatusResponse(lastForm))
}
