package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Voter struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name" validate:"required,min=2,max=100"`
	AdhaarNumber string             `bson:"adhaar_number" validate:"required,min=12"`
	Phone        string             `bson:"phone" validate:"regexp=(0|+91|091|91)[0-9]+$"`
	DateOfBirth  string             `bson:"date_of_birth" validate:"required"`
	Gender       string             `bson:"gender" validate:"required,eq=MALE|eq=FEMALE"`
	IsEligible   bool               `bson:"is_eligible" validate:"required"`
	IsVoted      bool               `bson:"is_voted" validate:"required"`
	CreatedAt    time.Time          `bson:"created_at" validate:"required"`
	UpdatedAt    time.Time          `bson:"updated_at" validate:"required"`
}
