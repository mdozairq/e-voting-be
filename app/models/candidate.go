package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Candidate struct {
	ID                       primitive.ObjectID `bson:"_id"`
	Email                    string             `bson:"email"`
	Username                 string             `bson:"username"`
	Phone                    string             `bson:"phone"`
	Password                 string             `bson:"password"`
	VoterID                  string             `bson:"voter_id"`
	PartyID                  string             `bson:"party_id"`
	RegisteredConstituencyID string             `bson:"registered_constituency_id"`
	Assets                   []string           `bson:"assets"`
	HasCrimeRecords          bool               `bson:"has_crime_records"`
	IsAccused                bool               `bson:"is_accused"`
	IsEligible               bool               `bson:"is_eligible"`
	IsRegistered             bool               `bson:"is_registered"`
	CreatedAt                time.Time          `bson:"created_at"`
	UpdatedAt                time.Time          `bson:"updated_at"`
}
