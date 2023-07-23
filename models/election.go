package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Election struct {
	ID             primitive.ObjectID `bson:"_id"`
	ElectionName   *string            `json:"election_name" validate:"required,min=2,max=100"`
	Description    *string            `json:"description" validate:"required,min=2,max=100"`
	ElectionType   *string            `json:"election_type" validate:"required"`
	ElectionStatus *string            `json:"election_status" validate:"required"`
	StartDate      *time.Time         `json:"start_date" validate:"required"`
	EndDate        *time.Time         `json:"end_date" validate:"required"`
	ElectionYear   *string            `json:"election_year" validate:"required"`
	IsRejected     *bool               `json:"is_rejected"`
	IsBypoll       *bool              `json:"is_voted" validate:"required"`
	CreatedAt      *time.Time         `json:"created_at" validate:"required"`
	UpdatedAt      *time.Time         `json:"updated_at" validate:"required"`
}
