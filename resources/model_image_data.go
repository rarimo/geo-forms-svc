/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type ImageData struct {
	// Image size. It cannot be more than 4 megabytes.
	ContentLength int64 `json:"content_length"`
	// Allowed content-type is `image/png`, `image/jpeg` or `image/x-jp2`
	ContentType string `json:"content_type"`
}
