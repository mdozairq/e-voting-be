package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/app/models"
	"github.com/mdozairq/e-voting-be/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ballotCollection = database.OpenCollection(database.Client, "ballot")


func CreateBallot() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			VoterID     string `json:"voter_id"`
			CandidateID string `json:"candidate_id"`
			PartyID     string `json:"party_id"`
			ElectionID  string `json:"election_id"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
			return
		}

		// Create a MongoDB context
		ctx := context.Background()

		objectID, err := primitive.ObjectIDFromHex(payload.VoterID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid party ID"})
			return
		}

		// Check if the voter exists in the voter database
		existingVoter := models.Voter{}
		err = voterCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingVoter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Voter not found"})
			return
		}

		// Check if the voter has already voted in this election
		existingBallot := models.Ballot{}
		err = ballotCollection.FindOne(ctx, bson.M{
			"voter_id":    payload.VoterID,
			"election_id": payload.ElectionID,
		}).Decode(&existingBallot)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Voter has already voted in this election"})
			return
		} else if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check voter's ballot"})
			return
		}

		// Create a new ballot
		ballot := models.Ballot{
			ID:          primitive.NewObjectID(),
			VoterID:     payload.VoterID,
			CandidateID: payload.CandidateID,
			PartyID:     payload.PartyID,
			ElectionID:  payload.ElectionID,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Insert the new ballot into the collection
		_, err = ballotCollection.InsertOne(ctx, ballot)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ballot"})
			return
		}

		// Update the voter's is_voted status
		_, err = voterCollection.UpdateOne(
			ctx,
			bson.M{"_id": objectID},
			bson.M{"$set": bson.M{"is_voted": true}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update voter's status"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Ballot created successfully"})
	}
}
