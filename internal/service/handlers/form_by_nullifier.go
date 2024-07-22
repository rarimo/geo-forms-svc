package handlers

import (
	"net/http"
	"time"

	"github.com/rarimo/geo-auth-svc/pkg/auth"
	"github.com/rarimo/geo-forms-svc/internal/data"
	"github.com/rarimo/geo-forms-svc/internal/service/requests"
	"github.com/rarimo/geo-forms-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func FormByNullifier(w http.ResponseWriter, r *http.Request) {
	nullifier, err := requests.NewFormByNullifier(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(nullifier)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	formStatus, err := FormsQ(r).Last(nullifier)
	if err != nil {
		Log(r).WithError(err).Error("Failed to get form by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if formStatus == nil {
		Log(r).Debugf("User %s doesn't have forms", nullifier)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	formStatus.NextFormAt = formStatus.CreatedAt.Add(Forms(r).Cooldown)

	ape.Render(w, newFormStatusResponse(formStatus))
}

func newFormStatusResponse(formStatus *data.FormStatus) resources.FormStatusResponse {
	untilNextForm := time.Now().UTC().Unix() - formStatus.NextFormAt.Unix()
	if untilNextForm < 0 {
		untilNextForm = 0
	}

	var processedAt *int64
	if formStatus.Status == data.ProcessedStatus {
		updatedAt := formStatus.UpdatedAt.Unix()
		processedAt = &updatedAt
	}

	return resources.FormStatusResponse{
		Data: resources.FormStatus{
			Key: resources.Key{
				ID:   formStatus.ID,
				Type: resources.FORM_STATUS,
			},
			Attributes: resources.FormStatusAttributes{
				Status:        formStatus.Status,
				CreatedAt:     formStatus.CreatedAt.Unix(),
				NextFormAt:    formStatus.NextFormAt.Unix(),
				UntilNextForm: untilNextForm,
				ProcessedAt:   processedAt,
			},
		},
	}
}
