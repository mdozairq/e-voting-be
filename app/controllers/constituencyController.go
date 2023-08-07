package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mdozairq/e-voting-be/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllConstituencies() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Assuming you have a MongoDB client and a collection reference named "constituencyCollection"
		// You need to pass your MongoDB client and collection reference as arguments to this function.

		// Create a MongoDB context
		ctx := context.Background()

		// Query all constituencies from the collection
		cursor, err := constituencyCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch constituencies"})
			return
		}
		defer cursor.Close(ctx)

		var constituencies []models.Constituency
		if err := cursor.All(ctx, &constituencies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode constituencies"})
			return
		}

		c.JSON(http.StatusOK, constituencies)
	}
}

func GetConstituencyByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the constituency ID from the request parameters
		constituencyID := c.Param("id")

		// Convert the constituency ID to an ObjectID
		objectID, err := primitive.ObjectIDFromHex(constituencyID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid constituency ID"})
			return
		}

		// Assuming you have a MongoDB client and a collection reference named "constituencyCollection"
		// You need to pass your MongoDB client and collection reference as arguments to this function.

		// Create a MongoDB context
		ctx := context.Background()

		// Query the constituency by its ID
		var constituency models.Constituency
		err = constituencyCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&constituency)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Constituency not found"})
			return
		}

		c.JSON(http.StatusOK, constituency)
	}
}
