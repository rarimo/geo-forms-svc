package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/rarimo/forms-svc/resources"
)

// 4 b64 letters encode 3 bytes, max image size = 12 MB -> (12/3)*4 * (1 << 20)
const maxImageSize = (1 << 20) * 16

func NewSubmitForm(r *http.Request) (req resources.SubmitFormRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":                validation.Validate(req.Data.Type, validation.Required, validation.In(resources.SUBMIT_FORM)),
		"data/attributes/name":     validation.Validate(req.Data.Attributes.Name, validation.Required, is.Alpha),
		"data/attributes/surname":  validation.Validate(req.Data.Attributes.Surname, validation.Required, is.Alpha),
		"data/attributes/id_num":   validation.Validate(req.Data.Attributes.IdNum, validation.Required, is.Digit),
		"data/attributes/birthday": validation.Validate(req.Data.Attributes.Birthday, validation.Required, validation.Date("01-02-2006")),
		"data/attributes/citizen":  validation.Validate(req.Data.Attributes.Citizen, validation.Required, is.Alpha),
		"data/attributes/visited":  validation.Validate(req.Data.Attributes.Visited, validation.Required),
		"data/attributes/purpose":  validation.Validate(req.Data.Attributes.Purpose, validation.Required),
		"data/attributes/country":  validation.Validate(req.Data.Attributes.Country, validation.Required, is.Alpha),
		"data/attributes/city":     validation.Validate(req.Data.Attributes.City, validation.Required, is.Alpha),
		"data/attributes/address":  validation.Validate(req.Data.Attributes.Address, validation.Required),
		"data/attributes/postal":   validation.Validate(req.Data.Attributes.Postal, validation.Required, is.Digit),
		"data/attributes/phone":    validation.Validate(req.Data.Attributes.Phone, validation.Required),
		"data/attributes/email":    validation.Validate(req.Data.Attributes.Email, validation.Required, is.Email),
		"data/attributes/image":    validation.Validate(req.Data.Attributes.Image, validation.Required, is.Base64, validation.Length(0, maxImageSize)),
	}

	return req, errs.Filter()
}

func newDecodeError(what string, err error) error {
	return validation.Errors{
		what: fmt.Errorf("decode request %s: %w", what, err),
	}
}
