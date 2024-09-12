/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type UploadImageResponseV2 struct {
	Key
	Attributes UploadImageResponseV2Attributes `json:"attributes"`
}
type UploadImageResponseV2Response struct {
	Data     UploadImageResponseV2 `json:"data"`
	Included Included              `json:"included"`
}

type UploadImageResponseV2ListResponse struct {
	Data     []UploadImageResponseV2 `json:"data"`
	Included Included                `json:"included"`
	Links    *Links                  `json:"links"`
	Meta     json.RawMessage         `json:"meta,omitempty"`
}

func (r *UploadImageResponseV2ListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *UploadImageResponseV2ListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustUploadImageResponseV2 - returns UploadImageResponseV2 from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUploadImageResponseV2(key Key) *UploadImageResponseV2 {
	var uploadImageResponseV2 UploadImageResponseV2
	if c.tryFindEntry(key, &uploadImageResponseV2) {
		return &uploadImageResponseV2
	}
	return nil
}
