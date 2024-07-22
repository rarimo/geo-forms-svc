package handlers

import (
	"net/http"
	"strings"

	"github.com/rarimo/geo-auth-svc/pkg/auth"
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

	lastStatus, err := FormsQ(r).Last(nullifier)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get last form")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if lastStatus == nil {
		Log(r).Debugf("Form for user=%s not found", nullifier)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if lastStatus.ID == id {
		ape.Render(w, newFormStatusResponse(lastStatus))
		return
	}

	formStatus, err := FormsQ(r).Get(id)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get form")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if formStatus == nil {
		Log(r).Debugf("Form with id=%s not found", id)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(formStatus.Nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}
	formStatus.NextFormAt = lastStatus.CreatedAt.Add(Forms(r).Cooldown)

	ape.Render(w, newFormStatusResponse(formStatus))
}
