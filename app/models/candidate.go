package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Candidate struct {
	ID                       primitive.ObjectID `bson:"_id" json:"_id"`
	Email                    string             `bson:"email" json:"email"`
	Username                 string             `bson:"username" json:"username"`
	Phone                    string             `bson:"phone" json:"phone"`
	Password                 string             `bson:"password" json:"password"`
	VoterID                  string             `bson:"voter_id" json:"voter_id"`
	PartyID                  string             `bson:"party_id" json:"party_id"`
	ElectionID               string             `bson:"election_id" json:"election_id"`
	RegisteredConstituencyID string             `bson:"registered_constituency_id" json:"registered_constituency_id"`
	Assets                   string             `bson:"assets" json:"assets"`
	HasCrimeRecords          bool               `bson:"has_crime_records" json:"has_crime_records"`
	IsAccused                bool               `bson:"is_accused" json:"is_accused"`
	IsEligible               bool               `bson:"is_eligible" json:"is_eligible"`
	IsRegistered             bool               `bson:"is_registered" json:"is_registered"`
	CreatedAt                time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt                time.Time          `bson:"updated_at" json:"updated_at"`
	Party                    Party              `bson:"-" json:"party"`
	Voter                    Voter              `bson:"-" json:"voter"`
}
