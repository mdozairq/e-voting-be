package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Party struct {
	ID             primitive.ObjectID `bson:"_id"`
	PartyName      *string            `json:"party_name" validate:"required,min=2,max=100"`
	PartyType      *string            `json:"party_type" validate:"required,min=2,max=100"`
	PartySlogan    *string            `json:"party_slogan" validate:"required`
	PartyLogo      *string            `json:"party_logo" validate:"required"`
	PartyLogoURL   *string            `json:"party_logo_url"`
	IsRuling       *bool              `json:"is_eligible" validate:"required`
	IsDisqualified *bool              `json:"is_voted" validate:"required"`
	CreatedAt      *time.Time         `json:"created_at" validate:"required"`
	UpdatedAt      *time.Time         `json:"updated_at" validate:"required"`
}
