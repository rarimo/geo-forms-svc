package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/rarimo/geo-forms-svc/internal/data"
	"github.com/rarimo/geo-forms-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func LastStatus(w http.ResponseWriter, r *http.Request) {
	nullifier := strings.ToLower(UserClaims(r)[0].Nullifier)

	form, err := FormsQ(r).FilterByNullifier(nullifier).Last()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get form by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if form == nil {
		Log(r).Debugf("User %s doesn't have forms", nullifier)
		ape.RenderErr(w, problems.NotFound())
		return
	}

	nextFormAt := form.CreatedAt
	if form.Status != data.CreatedStatus {
		nextFormAt = form.CreatedAt.Add(Forms(r).Cooldown)
	}

	ape.Render(w, newFormStatusResponse(*form, nextFormAt))
}

func newFormStatusResponse(form data.Form, nextFormAt time.Time) resources.FormResponse {
	untilNextForm := time.Now().UTC().Unix() - nextFormAt.Unix()
	if untilNextForm < 0 || form.Status == data.CreatedStatus {
		untilNextForm = 0
	}

	var processedAt *int64
	if form.Status == data.ProcessedStatus {
		updatedAt := form.UpdatedAt.Unix()
		processedAt = &updatedAt
	}

	return resources.FormResponse{
		Data: resources.Form{
			Key: resources.Key{
				ID:   form.ID,
				Type: resources.FORM,
			},
			Attributes: resources.FormAttributes{
				Address:       form.Address,
				Birthday:      form.Birthday,
				Citizen:       form.Citizen,
				City:          form.City,
				Country:       form.Country,
				Email:         form.Email,
				IdNum:         form.IDNum,
				Image:         form.Image,
				Name:          form.Name,
				Phone:         form.Phone,
				Postal:        form.Postal,
				Purpose:       form.Purpose,
				Surname:       form.Surname,
				Visited:       form.Visited,
				Status:        &form.Status,
				CreatedAt:     aws.Int64(form.CreatedAt.Unix()),
				NextFormAt:    aws.Int64(nextFormAt.Unix()),
				UntilNextForm: aws.Int64(untilNextForm),
				ProcessedAt:   processedAt,
			},
		},
	}
}
