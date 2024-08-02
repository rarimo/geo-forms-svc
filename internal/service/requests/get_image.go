package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gitlab.com/distributed_lab/urlval/v4"
)

type GetImage struct {
	ID     string
	ApiKey string `url:"api_key"`
}

func NewGetImage(r *http.Request) (req GetImage, err error) {
	req.ID = chi.URLParam(r, "id")

	if err = urlval.Decode(r.URL.Query(), &req); err != nil {
		err = newDecodeError("query", err)
		return
	}

	err = validation.Errors{
		"id": validation.Validate(req.ID, validation.Required, is.UUID),
	}.Filter()
	return

}
