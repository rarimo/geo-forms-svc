/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type UploadImageResponse struct {
	Key
	Attributes UploadImageResponseAttributes `json:"attributes"`
}
type UploadImageResponseResponse struct {
	Data     UploadImageResponse `json:"data"`
	Included Included            `json:"included"`
}

type UploadImageResponseListResponse struct {
	Data     []UploadImageResponse `json:"data"`
	Included Included              `json:"included"`
	Links    *Links                `json:"links"`
	Meta     json.RawMessage       `json:"meta,omitempty"`
}

func (r *UploadImageResponseListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *UploadImageResponseListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustUploadImageResponse - returns UploadImageResponse from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUploadImageResponse(key Key) *UploadImageResponse {
	var uploadImageResponse UploadImageResponse
	if c.tryFindEntry(key, &uploadImageResponse) {
		return &uploadImageResponse
	}
	return nil
}
