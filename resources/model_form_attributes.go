/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type FormAttributes struct {
	Address  string `json:"address"`
	Birthday string `json:"birthday"`
	Citizen  string `json:"citizen"`
	City     string `json:"city"`
	Country  string `json:"country"`
	// Form submission time. Unix time. Read-only.
	CreatedAt *int64 `json:"created_at,omitempty"`
	Email     string `json:"email"`
	IdNum     string `json:"id_num"`
	// base64 encoded image with max size 4 MB or URL for S3 storage with image up to 4 mb
	Image string `json:"image"`
	Name  string `json:"name"`
	// Time of the next possible form submission. Unix time. Read-only.
	NextFormAt *int64 `json:"next_form_at,omitempty"`
	// base64 encoded image with max size 4 MB or URL for S3 storage with image up to 4 mb
	PassportImage *string `json:"passport_image,omitempty"`
	Phone         string  `json:"phone"`
	Postal        string  `json:"postal"`
	// Form processing time. Absent if the status is accepted. Unix time. Read-only.
	ProcessedAt *int64 `json:"processed_at,omitempty"`
	Purpose     string `json:"purpose"`
	// Created - the empty form was created and now user can't use legacy submit Accepted - the data was saved by the service for further processing Processed - the data is processed and stored Read-only.
	Status  *string `json:"status,omitempty"`
	Surname string  `json:"surname"`
	// Time until the next form submission in seconds. Read-only.
	UntilNextForm *int64 `json:"until_next_form,omitempty"`
	Visited       string `json:"visited"`
}
