// models/aadhaar_model.go
package models

// AadhaarCard represents the data structure for an Aadhaar card.
type AdhaarCard struct {
	UID      string `bson:"uid" validate:"required,len=12"`
	Name     string `bson:"name" validate:"required"`
	DOB      string `bson:"dob" validate:"required,date"`
	Gender   string `bson:"gender" validate:"required,oneof=M F O"`
	City     string `bson:"city" validate:"required"`
	State    string `bson:"state" validate:"required"`
	Country  string `bson:"country" validate:"required"`
	Address  string `bson:"address" validate:"required"`
	PinCode  string `bson:"pincode" validate:"required"`
	PhotoURL string `json:"photo_url" bson:"photo_url" validate:"required"`
	Mobile   string `bson:"mobile" validate:"required,numeric,len=10"`
	Email    string `bson:"email,omitempty" validate:"omitempty,email"`
	Issuer   string `bson:"issuer" validate:"required"`
}
