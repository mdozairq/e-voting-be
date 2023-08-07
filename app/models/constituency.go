package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Constituency struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	City     string            `bson:"city" json:"city,omitempty"`
	District string            `bson:"district" json:"district,omitempty"`
	State    string            `bson:"state" json:"state,omitempty"`
	Country  string            `bson:"country" json:"country,omitempty"`
	PinCode  string            `bson:"pincode" json:"pincode,omitempty"`
}
