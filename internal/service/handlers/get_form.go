package handlers

import (
	"net/http"

	"github.com/rarimo/geo-auth-svc/pkg/auth"
	"github.com/rarimo/geo-forms-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetForm(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetForm(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	form, err := FormsQ(r).FilterByID(req.ID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get form")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if form == nil {
		Log(r).Debugf("Form with id=%s not found", req.ID)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(form.Nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	ape.Render(w, newFormResponse(form.ID))
}
