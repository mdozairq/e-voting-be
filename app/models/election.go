package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Election struct {
	ID            primitive.ObjectID `bson:"_id" json:"_id"`
	ElectionName  string             `bson:"election_name" json:"election_name" validate:"required,min=2,max=100"`
	Description   string             `bson:"description" json:"description" validate:"required,min=2,max=100"`
	ElectionType  string             `bson:"election_type" json:"election_type" validate:"required,oneof=GENERAL STATE MUNICIPAL PANCHAYAT"`
	Constituency  string             `bson:"constituency" json:"constituency" validate:"required"`
	StartDate     time.Time          `bson:"start_date" json:"start_date" validate:"required"`
	EndDate       time.Time          `bson:"end_date" json:"end_date" validate:"required"`
	ElectionYear  string             `bson:"election_year" json:"election_year" validate:"required"`
	IsActive      bool               `bson:"is_active" json:"is_active"`
	IsBypoll      bool               `json:"is_bypoll" json:"is_bypoll" validate:"required"`
	ElectionPhase string             `json:"election_phase" bson:"election_phase" validate:"required,oßneof=INITIALIZATION REGISTRATION VOTING RESULT DECLARED"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at" validate:"required"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at" validate:"required"`
}
