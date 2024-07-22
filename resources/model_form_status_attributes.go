/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type FormStatusAttributes struct {
	// Form submission time. Unix UTC time.
	CreatedAt int64 `json:"created_at"`
	// Time of the next possible form submission. Unix UTC time.
	NextFormAt int64 `json:"next_form_at"`
	// Form processing time. Absent if the status is accepted. Unix UTC time.
	ProcessedAt *int64 `json:"processed_at,omitempty"`
	// Accepted - the data was saved by the service for further processing Processed - the data is processed and stored
	Status string `json:"status"`
	// Time until the next form submission in seconds.
	UntilNextForm int64 `json:"until_next_form"`
}
