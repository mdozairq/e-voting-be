package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/app/models"
	"github.com/mdozairq/e-voting-be/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var partyCollection = database.OpenCollection(database.Client, "party")

func CreateParty() gin.HandlerFunc {
	return func(c *gin.Context) {
		var party models.Party
		if err := c.ShouldBindJSON(&party); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		party.ID = primitive.NewObjectID()
		party.CreatedAt = time.Now()
		party.UpdatedAt = time.Now()

		_, err := partyCollection.InsertOne(context.Background(), party)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create party"})
			return
		}

		c.JSON(http.StatusCreated, party)
	}
}

// func GetAllParties() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var parties []models.Party

// 		cursor, err := partyCollection.Find(context.Background(), bson.D{})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch parties"})
// 			return
// 		}
// 		defer cursor.Close(context.Background())

// 		for cursor.Next(context.Background()) {
// 			var party models.Party
// 			err := cursor.Decode(&party)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode parties"})
// 				return
// 			}
// 			parties = append(parties, party)
// 		}

// 		c.JSON(http.StatusOK, parties)
// 	}
// }

func GetAllParties() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the election ID from the query parameter
		electionID := c.Query("electionId")

		// Query to find candidates who are registered for the specified election
		matchCandidatesQuery := bson.D{
			{"election_id", electionID},
			{"is_registered", true},
		}

		// Find candidate documents based on the query
		var registeredCandidates []models.Candidate
		candidateCursor, err := candidateCollection.Find(context.Background(), matchCandidatesQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch registered candidates"})
			return
		}
		defer candidateCursor.Close(context.Background())

		// Collect the party IDs of registered candidates
		var partyIDs []primitive.ObjectID
		for candidateCursor.Next(context.Background()) {
			var candidate models.Candidate
			if err := candidateCursor.Decode(&candidate); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode candidate"})
				return
			}
			partyID, err := primitive.ObjectIDFromHex(candidate.PartyID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid party ID"})
				return
			}
			partyIDs = append(partyIDs, partyID)
			registeredCandidates = append(registeredCandidates, candidate)
		}
		log.Printf("Registered Candidates: %+v", partyIDs)

		// Query to find parties excluding those associated with registered candidates
		var matchPartiesQuery bson.D

		if len(partyIDs) > 0 {
			matchPartiesQuery = bson.D{
				{"_id", bson.D{
					{"$nin", partyIDs},
				}},
			}
		} else {
			// Handle the case where there are no partyIDs to exclude
			// For example, you might want to fetch all parties in this case
			matchPartiesQuery = bson.D{}
		}

		// Find party documents based on the query
		cursor, err := partyCollection.Find(context.Background(), matchPartiesQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch parties"})
			return
		}
		defer cursor.Close(context.Background())
		log.Printf("Not Registered Query: %+v", matchPartiesQuery)

		// Collect parties that match the query
		var parties []models.Party
		for cursor.Next(context.Background()) {
			var party models.Party
			if err := cursor.Decode(&party); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode parties"})
				return
			}
			parties = append(parties, party)
			log.Printf("Not Registered Party: %+v", party)
		}

		c.JSON(http.StatusOK, parties)
	}
}

func GetPartyByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		partyID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(partyID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid party ID"})
			return
		}

		var party models.Party
		err = partyCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&party)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Party not found"})
			return
		}

		c.JSON(http.StatusOK, party)
	}
}

func DeleteParty() gin.HandlerFunc {
	return func(c *gin.Context) {
		partyID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(partyID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid party ID"})
			return
		}

		_, err = partyCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete party"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Party deleted"})
	}
}
