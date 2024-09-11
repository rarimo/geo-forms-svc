package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/urlval/v4"
)

func NewGetImage(r *http.Request) (id, apiKey string, err error) {
	id = chi.URLParam(r, "id")

	var queryParams struct {
		APIKey string `url:"api"`
	}
	if err = urlval.Decode(r.URL.Query(), &queryParams); err != nil {
		err = newDecodeError("query", err)
		return
	}
	apiKey = queryParams.APIKey
	return
}
