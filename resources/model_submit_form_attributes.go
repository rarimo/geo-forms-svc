/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type SubmitFormAttributes struct {
	Address string `json:"address"`
	// Date formated as DD/MM/YYYY
	Birthday string `json:"birthday"`
	Citizen  string `json:"citizen"`
	City     string `json:"city"`
	Country  string `json:"country"`
	Email    string `json:"email"`
	IdNum    string `json:"id_num"`
	// For default endpoint:   base64 encoded image with max size 4 MB; For lightweight endpoint:   link to the image in s3 storage;
	Image   string `json:"image"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Postal  string `json:"postal"`
	Purpose string `json:"purpose"`
	Surname string `json:"surname"`
	// Date formated as DD/MM/YYYY
	Visited string `json:"visited"`
}
