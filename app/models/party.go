package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Party struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id"`
	PartyName      string            `bson:"party_name" json:"party_name" validate:"required,min=2,max=100"`
	PartyType      string            `bson:"party_type" json:"party_type" validate:"required,min=2,max=100"`
	PartySlogan    string            `bson:"party_slogan" json:"party_slogan" validate:"required"`
	PartyLogo      string            `bson:"party_logo" json:"party_logo" validate:"required"`
	PartyLogoURL   string            `bson:"party_logo_url" json:"party_logo_url"`
	IsRuling       bool              `bson:"is_ruling" json:"is_ruling" validate:"required"`
	IsDisqualified bool              `bson:"is_disqualified" json:"is_disqualified" validate:"required"`
	CreatedAt      time.Time         `bson:"created_at" json:"created_at" validate:"required"`
	UpdatedAt      time.Time         `bson:"updated_at" json:"updated_at" validate:"required"`
}
