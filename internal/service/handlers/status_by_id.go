package handlers

import (
	"net/http"
	"strings"

	"github.com/rarimo/geo-auth-svc/pkg/auth"
	"github.com/rarimo/geo-forms-svc/internal/data"
	"github.com/rarimo/geo-forms-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func StatusByID(w http.ResponseWriter, r *http.Request) {
	id, err := requests.NewStatusByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	nullifier := strings.ToLower(UserClaims(r)[0].Nullifier)

	lastForm, err := FormsQ(r).FilterByNullifier(nullifier).Last()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get last form")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if lastForm == nil {
		Log(r).Debugf("Form for user=%s not found", nullifier)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	nextFormAt := lastForm.CreatedAt
	if lastForm.Status != data.CreatedStatus {
		nextFormAt = lastForm.CreatedAt.Add(Forms(r).Cooldown)
	}

	if lastForm.ID == id {
		ape.Render(w, newFormStatusResponse(*lastForm, nextFormAt))
		return
	}

	form, err := FormsQ(r).FilterByID(id).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get form")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if form == nil {
		Log(r).Debugf("Form with id=%s not found", id)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(form.Nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	ape.Render(w, newFormStatusResponse(*form, nextFormAt))
}
