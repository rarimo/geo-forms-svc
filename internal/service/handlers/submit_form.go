package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/rarimo/geo-forms-svc/internal/data"
	"github.com/rarimo/geo-forms-svc/internal/service/requests"
	"github.com/rarimo/geo-forms-svc/resources"
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

	formStatus, err := FormsQ(r).Last(nullifier)
	if err != nil {
		Log(r).WithError(err).Errorf("Failed to get last user form for nullifier [%s]", nullifier)
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if formStatus != nil {
		next := formStatus.CreatedAt.Add(Forms(r).Cooldown)
		if next.After(time.Now().UTC()) {
			Log(r).Debugf("Form submitted time: %s; next available time: %s", formStatus.CreatedAt, next)
			ape.RenderErr(w, problems.TooManyRequests())
			return
		}
	}

	userData := req.Data.Attributes
	form := &data.Form{
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
		Image:     &userData.Image,
	}

	if err = Forms(r).SendForms(form); err != nil {
		Log(r).WithError(err).Error("Failed to send form")
		form.Status = data.AcceptedStatus
	}

	_, err = FormsQ(r).Insert(form)
	if err != nil {
		Log(r).WithError(err).Error("Failed to insert form")
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
				Status:   &form.Status,
				Address:  form.Address,
				Birthday: form.Birthday,
				Citizen:  form.Citizen,
				City:     form.City,
				Country:  form.Country,
				Email:    form.Email,
				IdNum:    form.IDNum,
				Name:     form.Name,
				Phone:    form.Phone,
				Postal:   form.Postal,
				Purpose:  form.Purpose,
				Surname:  form.Surname,
				Visited:  form.Visited,
			},
		},
	}
}
