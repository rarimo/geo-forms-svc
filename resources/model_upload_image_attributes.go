/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UploadImageAttributes struct {
	// Image size. It cannot be more than 4 megabytes.
	ContentLength int64 `json:"content_length"`
	// Allowed content-type is `image/png` or `image/jpeg`
	ContentType string `json:"content_type"`
}
