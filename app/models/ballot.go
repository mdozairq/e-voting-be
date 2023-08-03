package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Ballot struct {
	ID          primitive.ObjectID `bson:"_id"`
	VoterID     string             `json:"voter_id"`
	CandidateID string             `json:"candidate_id"`
	PartyID     string             `json:"party_id"`
	ElectionID  string             `json:"election_id"`
	IsEligible  bool               `json:"is_eligible"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
