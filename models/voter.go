package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Voter struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         *string            `json:"name" validate:"required,min=2,max=100"`
	AdhaarNumber *string            `json:"adhaar_number" validate:"required,min=12`
	Phone        *string            `json:"phone" validate:"regexp=(0|\+91|091|91)[0-9]+$"`
	DateOfBirth  *time.Time         `json:"date_of_birth" validate:"required"`
	Gender       *string            `json:"gender" validate:"required,eq=MALE|eq=FEMALE"`
	IsEligible   *bool              `json:"is_eligible" validate:"required`
	IsVoted      *bool              `json:"is_voted" validate:"required"`
	CreatedAt    time.Time         `json:"created_at" validate:"required"`
	UpdatedAt    time.Time         `json:"updated_at" validate:"required"`
}
