package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Constituency struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	City     *string            `json:"city"`
	District *string            `json:"district"`
	State    *string            `json:"state"`
	Country  *string            `json:"country"`
	PinCode  *string            `json:"pincode"`
}
