package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Voter struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	Name         string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	AdhaarNumber string             `bson:"adhaar_number" json:"adhaar_number" validate:"required,min=12"`
	Phone        string             `bson:"phone" json:"phone" validate:"regexp=(0|+91|091|91)[0-9]+$"`
	DateOfBirth  string             `bson:"date_of_birth" json:"date_of_birth" validate:"required"`
	Gender       string             `bson:"gender" json:"gender" validate:"required,eq=MALE|eq=FEMALE"`
	IsEligible   bool               `bson:"is_eligible" json:"is_eligible" validate:"required"`
	IsVoted      bool               `bson:"is_voted" json:"is_voted" validate:"required"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at" validate:"required"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at" validate:"required"`
}
