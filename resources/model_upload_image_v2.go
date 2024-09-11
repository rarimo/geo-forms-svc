/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type UploadImageV2 struct {
	Key
	Attributes UploadImageV2Attributes `json:"attributes"`
}
type UploadImageV2Request struct {
	Data     UploadImageV2 `json:"data"`
	Included Included      `json:"included"`
}

type UploadImageV2ListRequest struct {
	Data     []UploadImageV2 `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *UploadImageV2ListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *UploadImageV2ListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustUploadImageV2 - returns UploadImageV2 from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUploadImageV2(key Key) *UploadImageV2 {
	var uploadImageV2 UploadImageV2
	if c.tryFindEntry(key, &uploadImageV2) {
		return &uploadImageV2
	}
	return nil
}
