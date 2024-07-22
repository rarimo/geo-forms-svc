/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type UploadImage struct {
	Key
	Attributes UploadImageAttributes `json:"attributes"`
}
type UploadImageRequest struct {
	Data     UploadImage `json:"data"`
	Included Included    `json:"included"`
}

type UploadImageListRequest struct {
	Data     []UploadImage   `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *UploadImageListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *UploadImageListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustUploadImage - returns UploadImage from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUploadImage(key Key) *UploadImage {
	var uploadImage UploadImage
	if c.tryFindEntry(key, &uploadImage) {
		return &uploadImage
	}
	return nil
}
