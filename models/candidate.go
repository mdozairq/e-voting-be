package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Candidate struct {
	ID                       primitive.ObjectID `bson:"_id"`
	VoterID                  string             `json:"voter_id"`
	PartyID                  string             `json:"party_id"`
	RegisteredConstituencyID string             `json:"registered_constituency_id"`
	Assets                   []string           `json:"assets"`
	HasCrimeRecords          bool               `json:"has_crime_records"`
	IsAccused                bool               `json:"is_accused"`
	IsEligible               bool               `json:"is_eligible"`
	CreatedAt                time.Time          `json:"created_at"`
	UpdatedAt                time.Time          `json:"updated_at"`
}
