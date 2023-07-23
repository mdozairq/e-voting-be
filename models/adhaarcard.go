// models/aadhaar_model.go
package models

// AadhaarCard represents the data structure for an Aadhaar card.
type AdhaarCard struct {
	UID      string `json:"_id" validate:"required,len=12"`
	Name     string `json:"name" validate:"required"`
	DOB      string `json:"dob" validate:"required,date"`
	Gender   string `json:"gender" validate:"required,oneof=M F O"`
	City     string `json:"city" validate:"required"`
	State    string `json:"state" validate:"required"`
	Country  string `json:"country" validate:"required"`
	Street   string `json:"street" validate:"required"`
	PhotoURL string `json:"photo_url" validate:"required,url"`
	Mobile   string `json:"mobile" validate:"required,numeric,len=10"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Issuer   string `json:"issuer" validate:"required"`
}
