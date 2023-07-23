package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OTP represents the structure of the OTP document.
type OTP struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UID         string             `bson:"uid"`
	PhoneNumber string             `bson:"phone_number"`
	OTP         string             `bson:"otp"`
	Expiration  time.Time          `bson:"expiration"`
	CreatedAt   time.Time          `bson:"created_at"`
}
