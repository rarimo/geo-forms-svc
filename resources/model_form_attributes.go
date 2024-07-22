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
	Email    string `json:"email"`
	IdNum    string `json:"id_num"`
	// base64 encoded image with max size 4 MB or URL for S3 storage with image up to 4 mb
	Image   string `json:"image"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Postal  string `json:"postal"`
	Purpose string `json:"purpose"`
	// Accepted - the data was saved by the service for further processing Processed - the data is processed and stored Read-only.
	Status  *string `json:"status,omitempty"`
	Surname string  `json:"surname"`
	Visited string  `json:"visited"`
}
