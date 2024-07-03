package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type GetForm struct {
	ID string
}

func NewGetForm(r *http.Request) (req GetForm, err error) {
	req.ID = chi.URLParam(r, "id")

	err = validation.Errors{
		"id": validation.Validate(req.ID, validation.Required, is.UUID),
	}.Filter()
	return
}
