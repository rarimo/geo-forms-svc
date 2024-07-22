package handlers

import (
	"net/http"

	"github.com/rarimo/geo-auth-svc/pkg/auth"
	"github.com/rarimo/geo-forms-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func FormByID(w http.ResponseWriter, r *http.Request) {
	id, err := requests.NewFormByID(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
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

	// formStatusByNullifier will never be nil because of the previous logic
	lastFormStatus, err := FormsQ(r).Last(formStatus.Nullifier)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get last form")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	formStatus.NextFormAt = lastFormStatus.CreatedAt.Add(Forms(r).Cooldown)

	ape.Render(w, newFormStatusResponse(formStatus))
}
