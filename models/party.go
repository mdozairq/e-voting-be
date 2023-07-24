package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Party struct {
	ID             primitive.ObjectID `bson:"_id"`
	PartyName      *string            `bson:"party_name" validate:"required,min=2,max=100"`
	PartyType      *string            `bson:"party_type" validate:"required,min=2,max=100"`
	PartySlogan    *string            `bson:"party_slogan" validate:"required`
	PartyLogo      *string            `bson:"party_logo" validate:"required"`
	PartyLogoURL   *string            `bson:"party_logo_url"`
	IsRuling       *bool              `bson:"is_eligible" validate:"required`
	IsDisqualified *bool              `bson:"is_voted" validate:"required"`
	CreatedAt      *time.Time         `bson:"created_at" validate:"required"`
	UpdatedAt      *time.Time         `bson:"updated_at" validate:"required"`
}
