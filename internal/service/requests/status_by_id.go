package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func NewStatusByID(r *http.Request) (id string, err error) {
	id = chi.URLParam(r, "id")

	err = validation.Errors{
		"id": validation.Validate(id, validation.Required, is.UUID),
	}.Filter()
	return
}
