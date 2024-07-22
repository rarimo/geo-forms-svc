/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type FormStatus struct {
	Key
	Attributes FormStatusAttributes `json:"attributes"`
}
type FormStatusResponse struct {
	Data     FormStatus `json:"data"`
	Included Included   `json:"included"`
}

type FormStatusListResponse struct {
	Data     []FormStatus    `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *FormStatusListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *FormStatusListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustFormStatus - returns FormStatus from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustFormStatus(key Key) *FormStatus {
	var formStatus FormStatus
	if c.tryFindEntry(key, &formStatus) {
		return &formStatus
	}
	return nil
}
