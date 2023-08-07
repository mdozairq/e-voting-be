package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mdozairq/e-voting-be/app/models"
	"github.com/mdozairq/e-voting-be/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var electionCollection = database.OpenCollection(database.Client, "elections")

type InitializeElectionDto struct {
	ElectionName string    `json:"election_name" validate:"required,min=2,max=100"`
	Description  string    `json:"description" validate:"required,min=2,max=100"`
	ElectionType string    `json:"election_type" validate:"required,oneof=GENERAL STATE MUNICIPAL PANCHAYAT"`
	StartDate    time.Time `json:"start_date" validate:"required"`
	EndDate      time.Time `json:"end_date" validate:"required"`
	ElectionYear string    `json:"election_year" validate:"required"`
	IsActive     bool      `json:"is_active"`
	IsBypoll     bool      `json:"is_bypoll"`
}

func InitializeElection() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var input InitializeElectionDto
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"binding error": err.Error()})
			return
		}

		fmt.Println(input)
		var validate = validator.New()
		if err := validate.Struct(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"validation error": err.Error()})
			return
		}

		existingElection := models.Election{}
		err := electionCollection.FindOne(ctx, bson.M{"election_year": input.ElectionYear, "election_name": input.ElectionName, "election_type": input.ElectionType}).Decode(&existingElection)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "An election with the same year and name type already exists"})
			return
		}

		election := models.Election{
			ID:            primitive.NewObjectID(),
			ElectionName:  input.ElectionName,
			Description:   input.Description,
			ElectionType:  input.ElectionType,
			StartDate:     input.StartDate,
			EndDate:       input.EndDate,
			ElectionYear:  input.ElectionYear,
			IsBypoll:      input.IsBypoll,
			ElectionPhase: "INITIALIZATION",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		_, err = electionCollection.InsertOne(ctx, election)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create election"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Election initialized successfully"})
	}
}

func GetAllElections() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the page number and page size from query parameters
		pageNumber, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || pageNumber < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		if err != nil || pageSize < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
			return
		}

		// Calculate the skip count based on page number and page size
		skip := (pageNumber - 1) * pageSize

		// Assuming you have a function to fetch all elections from the database
		allElections, total, err := findAllElections(skip, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch elections"})
			return
		}

		// Create the response data with elections and total count
		responseData := gin.H{
			"total":     total,
			"page":      pageNumber,
			"page_size": pageSize,
			"elections": allElections,
		}

		// Return the response data as JSON
		c.JSON(http.StatusOK, responseData)
	}
}

func findAllElections(skip, pageSize int) ([]models.Election, int64, error) {
	var allElections []models.Election
	var total int64

	// Assuming you have a MongoDB client and a collection reference named "electionCollection"
	// You need to pass your MongoDB client and collection reference as arguments to this function.

	// Create a MongoDB context
	ctx := context.Background()

	// Define options for pagination
	options := options.Find()
	options.SetSkip(int64(skip))
	options.SetLimit(int64(pageSize))

	// Fetch all elections with pagination
	cursor, err := electionCollection.Find(ctx, bson.D{}, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Decode results into a slice of Election
	for cursor.Next(ctx) {
		var election models.Election
		err := cursor.Decode(&election)
		if err != nil {
			return nil, 0, err
		}
		allElections = append(allElections, election)
	}

	// Get the total count of elections
	total, err = electionCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	return allElections, total, nil
}

func GetElectionByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the election ID from the request parameters
		electionID := c.Param("id")

		// Convert the election ID to an ObjectID
		objectID, err := primitive.ObjectIDFromHex(electionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID"})
			return
		}

		// Assuming you have a MongoDB client and collection references named "electionCollection" and "constituencyCollection"
		// You need to pass your MongoDB client and collection references as arguments to this function.

		// Create a MongoDB context
		ctx := context.Background()

		// Find the election based on its ID
		var election models.Election
		err = electionCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&election)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
			return
		}

		// If the election has a Constituency ID, fetch the constituency details
		if election.Constituency != "" {
			// Convert the Constituency ID to an ObjectID
			constituencyID, err := primitive.ObjectIDFromHex(election.Constituency)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch constituency details"})
				return
			}

			// Find the constituency based on its ID
			var constituency models.Constituency
			err = constituencyCollection.FindOne(ctx, bson.M{"_id": constituencyID}).Decode(&constituency)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch constituency details"})
				return
			}

			// Set the Constituency details in the Election model
			electionDetails := struct {
				ID            primitive.ObjectID `bson:"_id" json:"_id"`
				ElectionName  string             `bson:"election_name" json:"election_name" validate:"required,min=2,max=100"`
				Description   string             `bson:"description" json:"description" validate:"required,min=2,max=100"`
				ElectionType  string             `bson:"election_type" json:"election_type" validate:"required,oneof=GENERAL STATE MUNICIPAL PANCHAYAT"`
				Constituency  models.Constituency `bson:"-" json:"constituency" validate:"required"`
				StartDate     time.Time          `bson:"start_date" json:"start_date" validate:"required"`
				EndDate       time.Time          `bson:"end_date" json:"end_date" validate:"required"`
				ElectionYear  string             `bson:"election_year" json:"election_year" validate:"required"`
				IsActive      bool               `bson:"is_active" json:"is_active"`
				IsBypoll      bool               `json:"is_bypoll" json:"is_bypoll" validate:"required"`
				ElectionPhase string             `json:"election_phase" bson:"election_phase" validate:"required,oneof=INITIALIZATION REGISTRATION VOTING RESULT"`
				CreatedAt     time.Time          `bson:"created_at" json:"created_at" validate:"required"`
				UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at" validate:"required"`
			}{
				ID:            election.ID,
				ElectionName:  election.ElectionName,
				Description:   election.Description,
				ElectionType:  election.ElectionType,
				Constituency:  constituency,
				StartDate:     election.StartDate,
				EndDate:       election.EndDate,
				ElectionYear:  election.ElectionYear,
				IsActive:      election.IsActive,
				IsBypoll:      election.IsBypoll,
				ElectionPhase: election.ElectionPhase,
				CreatedAt:     election.CreatedAt,
				UpdatedAt:     election.UpdatedAt,
			}

			c.JSON(http.StatusOK, electionDetails)
		} else {
			c.JSON(http.StatusOK, election)
		}
	}
}

