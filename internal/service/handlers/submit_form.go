package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/rarimo/forms-svc/internal/data"
	"github.com/rarimo/forms-svc/internal/service/requests"
	"github.com/rarimo/forms-svc/resources"
	"github.com/rarimo/geo-auth-svc/pkg/auth"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func SubmitForm(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewSubmitForm(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	nullifier := strings.ToLower(req.Data.ID)
	if !auth.Authenticates(UserClaims(r), auth.VerifiedGrant(req.Data.ID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	form, err := FormsQ(r).FilterByNullifier(nullifier).Last()
	if err != nil {
		Log(r).WithError(err).Errorf("failed to get last user form for nullifier [%s]", nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if form != nil && form.CreatedAt.Add(Forms(r).Cooldown).After(time.Now().UTC()) {
		Log(r).Debugf("Form submitted time: %s; Next available time: %s",
			form.CreatedAt.String(),
			form.CreatedAt.Add(Forms(r).Cooldown).String())
		ape.RenderErr(w, problems.TooManyRequests())
		return
	}

	userData := req.Data.Attributes
	form = &data.Form{
		Nullifier: nullifier,
		Status:    data.ProcessedStatus,
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
		Image:     userData.Image,
	}

	err = Forms(r).SendForms(*form)
	if err != nil {
		Log(r).WithError(err).Error("failed to send form")
		form.Status = data.AcceptedStatus
	}

	form, err = FormsQ(r).Insert(*form)
	if err != nil {
		Log(r).WithError(err).Error("failed to insert form")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newFormResponse(form))
}

func newFormResponse(form *data.Form) resources.FormResponse {
	return resources.FormResponse{
		Data: resources.Form{
			Key: resources.Key{
				ID:   form.ID,
				Type: resources.FORM,
			},
			Attributes: resources.FormAttributes{
				Name:     form.Name,
				Surname:  form.Surname,
				IdNum:    form.IDNum,
				Birthday: form.Birthday,
				Citizen:  form.Citizen,
				Visited:  form.Visited,
				Purpose:  form.Purpose,
				Country:  form.Country,
				City:     form.City,
				Address:  form.Address,
				Postal:   form.Postal,
				Phone:    form.Phone,
				Email:    form.Email,
				Image:    form.Image,
				Status:   form.Status,
			},
		},
	}
}
