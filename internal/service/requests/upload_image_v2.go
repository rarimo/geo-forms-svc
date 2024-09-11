package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/geo-forms-svc/resources"
)

func NewUploadImageV2(r *http.Request) (req resources.UploadImageV2Request, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	imageData := req.Data.Attributes.SelfieImage
	errs := validation.Errors{
		"data/type": validation.Validate(req.Data.Type, validation.Required, validation.In(resources.UPLOAD_IMAGE)),
		"data/attributes/selfie_image/content_type":   validation.Validate(imageData.ContentType, validation.Required, validation.In("image/png", "image/jpeg", "image/x-jp2")),
		"data/attributes/selfie_image/content_length": validation.Validate(imageData.ContentLength, validation.Required, validation.Min(int64(1)), validation.Max(maxImageSize)),
	}

	if req.Data.Attributes.PassportImage != nil {
		imageData = *req.Data.Attributes.PassportImage
		errs["data/attributes/passport_image/content_type"] = validation.Validate(imageData.ContentType, validation.Required, validation.In("image/png", "image/jpeg", "image/x-jp2"))
		errs["data/attributes/passport_image/content_length"] = validation.Validate(imageData.ContentLength, validation.Required, validation.Min(int64(1)), validation.Max(maxImageSize))
	}

	return req, errs.Filter()
}
