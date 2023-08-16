package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ballot struct {
	ID          primitive.ObjectID `bson:"_id"`
	VoterID     string             `bson:"voter_id" json:"voter_id"`
	CandidateID string             `bson:"candidate_id" json:"candidate_id" `
	PartyID     string             `bson:"party_id" json:"party_id"`
	ElectionID  string             `bson:"election_id" json:"election_id"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
