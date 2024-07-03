/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type SubmitForm struct {
	Key
	Attributes SubmitFormAttributes `json:"attributes"`
}
type SubmitFormRequest struct {
	Data     SubmitForm `json:"data"`
	Included Included   `json:"included"`
}

type SubmitFormListRequest struct {
	Data     []SubmitForm    `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *SubmitFormListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *SubmitFormListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustSubmitForm - returns SubmitForm from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSubmitForm(key Key) *SubmitForm {
	var submitForm SubmitForm
	if c.tryFindEntry(key, &submitForm) {
		return &submitForm
	}
	return nil
}
