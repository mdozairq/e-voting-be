package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Election struct {
	ID             primitive.ObjectID `bson:"_id"`
	ElectionName   *string            `bson:"election_name" validate:"required,min=2,max=100"`
	Description    *string            `bson:"description" validate:"required,min=2,max=100"`
	ElectionType   *string            `bson:"election_type" validate:"required"`
	ElectionStatus *string            `bson:"election_status" validate:"required"`
	StartDate      *time.Time         `bson:"start_date" validate:"required"`
	EndDate        *time.Time         `bson:"end_date" validate:"required"`
	ElectionYear   *string            `bson:"election_year" validate:"required"`
	IsRejected     *bool               `bson:"is_rejected"`
	IsBypoll       *bool              `bson:"is_voted" validate:"required"`
	CreatedAt      *time.Time         `bson:"created_at" validate:"required"`
	UpdatedAt      *time.Time         `bson:"updated_at" validate:"required"`
}
