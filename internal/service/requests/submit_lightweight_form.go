package requests

import (
	"encoding/json"
	"net/http"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rarimo/geo-forms-svc/resources"
)

func NewSubmitLightweightForm(r *http.Request) (req resources.SubmitFormRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":                validation.Validate(req.Data.Type, validation.Required, validation.In(resources.SUBMIT_FORM)),
		"data/attributes/name":     validation.Validate(req.Data.Attributes.Name, validation.Required),
		"data/attributes/surname":  validation.Validate(req.Data.Attributes.Surname, validation.Required),
		"data/attributes/id_num":   validation.Validate(req.Data.Attributes.IdNum, validation.Required),
		"data/attributes/birthday": validation.Validate(req.Data.Attributes.Birthday, validation.Required),
		"data/attributes/citizen":  validation.Validate(req.Data.Attributes.Citizen, validation.Required),
		"data/attributes/visited":  validation.Validate(req.Data.Attributes.Visited, validation.Required),
		"data/attributes/purpose":  validation.Validate(req.Data.Attributes.Purpose, validation.Required),
		"data/attributes/country":  validation.Validate(req.Data.Attributes.Country, validation.Required),
		"data/attributes/city":     validation.Validate(req.Data.Attributes.City, validation.Required),
		"data/attributes/address":  validation.Validate(req.Data.Attributes.Address, validation.Required),
		"data/attributes/postal":   validation.Validate(req.Data.Attributes.Postal, validation.Required),
		"data/attributes/phone":    validation.Validate(req.Data.Attributes.Phone, validation.Required),
		"data/attributes/email":    validation.Validate(req.Data.Attributes.Email, validation.Required, validation.Match(regexp.MustCompile(`[A-Z0-9a-z._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,64}`))),
		"data/attributes/image":    validation.Validate(req.Data.Attributes.Image, validation.Required, is.URL),
	}

	return req, errs.Filter()
}
