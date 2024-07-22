package requests

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var NullifierRegexp = regexp.MustCompile("^0x[0-9a-fA-F]{64}$")

func NewFormByNullifier(r *http.Request) (nullifier string, err error) {
	nullifier = chi.URLParam(r, "nullifier")

	err = validation.Errors{
		"nullifier": validation.Validate(nullifier, validation.Required, validation.Match(NullifierRegexp)),
	}.Filter()
	return
}
