/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Form struct {
	Key
	Attributes FormAttributes `json:"attributes"`
}
type FormResponse struct {
	Data     Form     `json:"data"`
	Included Included `json:"included"`
}

type FormListResponse struct {
	Data     []Form          `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *FormListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *FormListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustForm - returns Form from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustForm(key Key) *Form {
	var form Form
	if c.tryFindEntry(key, &form) {
		return &form
	}
	return nil
}
