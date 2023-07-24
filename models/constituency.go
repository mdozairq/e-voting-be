package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Constituency struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	City     *string            `bson:"city"`
	District *string            `bson:"district"`
	State    *string            `bson:"state"`
	Country  *string            `bson:"country"`
	PinCode  *string            `bson:"pincode"`
}
