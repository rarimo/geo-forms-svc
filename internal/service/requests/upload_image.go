package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/geo-forms-svc/resources"
)

const maxImageSize = int64(1 << 22)

func NewUploadImage(r *http.Request) (req resources.UploadImageRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":                      validation.Validate(req.Data.Type, validation.Required, validation.In(resources.UPLOAD_IMAGE)),
		"data/attributes/content_type":   validation.Validate(req.Data.Attributes.ContentType, validation.Required, validation.In("image/png", "image/jpeg")),
		"data/attributes/content_length": validation.Validate(req.Data.Attributes.ContentLength, validation.Required, validation.Min(int64(1)), validation.Max(maxImageSize)),
	}

	return req, errs.Filter()
}
