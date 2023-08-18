package controllers

import (
	"context"
	"fmt"
	"log"
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
	Constituency string    `bson:"constituency" json:"constituency" validate:"required"`
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
			Constituency:  input.Constituency,
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

		c.JSON(http.StatusCreated, gin.H{"data": &election, "message": "Election initialized successfully"})
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
				ID            primitive.ObjectID  `bson:"_id" json:"_id"`
				ElectionName  string              `bson:"election_name" json:"election_name" validate:"required,min=2,max=100"`
				Description   string              `bson:"description" json:"description" validate:"required,min=2,max=100"`
				ElectionType  string              `bson:"election_type" json:"election_type" validate:"required,oneof=GENERAL STATE MUNICIPAL PANCHAYAT"`
				Constituency  models.Constituency `bson:"-" json:"constituency" validate:"required"`
				StartDate     time.Time           `bson:"start_date" json:"start_date" validate:"required"`
				EndDate       time.Time           `bson:"end_date" json:"end_date" validate:"required"`
				ElectionYear  string              `bson:"election_year" json:"election_year" validate:"required"`
				IsActive      bool                `bson:"is_active" json:"is_active"`
				IsBypoll      bool                `json:"is_bypoll" json:"is_bypoll" validate:"required"`
				ElectionPhase string              `json:"election_phase" bson:"election_phase" validate:"required,oneof=INITIALIZATION REGISTRATION VOTING RESULT"`
				CreatedAt     time.Time           `bson:"created_at" json:"created_at" validate:"required"`
				UpdatedAt     time.Time           `bson:"updated_at" json:"updated_at" validate:"required"`
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

func GetRegistrationElections() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Assuming you have a MongoDB client and a collection reference named "electionCollection"
		// You need to pass your MongoDB client and collection reference as arguments to this function.

		// Create a MongoDB context
		ctx := context.Background()

		// Query elections with the "REGISTRATION" phase
		filter := bson.D{{"election_phase", "REGISTRATION"}}
		cursor, err := electionCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch elections"})
			return
		}
		defer cursor.Close(ctx)

		// Decode the results into a slice of Election
		var registrationElections []models.Election
		for cursor.Next(ctx) {
			var election models.Election
			err := cursor.Decode(&election)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode results"})
				return
			}
			registrationElections = append(registrationElections, election)
		}

		c.JSON(http.StatusOK, gin.H{"elections": registrationElections})
	}
}

func UpdateElectionPhase() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the election ID from the request parameters
		electionID := c.Param("id")
		log.Printf("ElectionId: %+v", electionID)
		// Parse the request body to get the start and end time
		var updatePhase struct {
			StartTime time.Time `json:"start_data" validate:"required"`
			EndTime   time.Time `json:"end_date" validate:"required"`
		}
		if err := c.ShouldBindJSON(&updatePhase); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Convert the election ID to an ObjectID
		objectID, err := primitive.ObjectIDFromHex(electionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID"})
			return
		}

		// Create a MongoDB context
		ctx := context.Background()

		// Find the existing election to check the current phase
		var existingElection models.Election
		err = electionCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingElection)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
			return
		}

		// Define the allowed phase sequence
		allowedPhases := []string{"INITIALIZATION", "REGISTRATION", "VOTING", "RESULT", "DECLARED"}

		// Find the index of the current phase in the sequence
		currentPhaseIndex := -1
		for i, phase := range allowedPhases {
			if phase == existingElection.ElectionPhase {
				currentPhaseIndex = i
				break
			}
		}

		// If the current phase is not found or it's the last phase (RESULT), return "Result declared"
		if currentPhaseIndex == -1 || currentPhaseIndex == len(allowedPhases)-1 {
			c.JSON(http.StatusOK, gin.H{"message": "Result declared"})
			return
		}

		// Update the election phase to the next phase in the sequence
		nextPhase := allowedPhases[currentPhaseIndex+1]

		// Update the election phase, start time, and end time
		update := bson.D{
			{"$set", bson.D{
				{"election_phase", nextPhase},
				{"start_date", updatePhase.StartTime},
				{"end_date", updatePhase.EndTime},
			}},
		}
		_, err = electionCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update election phase"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Election phase updated to %s", nextPhase)})
	}
}

func GetElectionResults() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the election ID from the query parameter
		electionID := c.Param("electionId")
		log.Printf("ElectionId: %+v", electionID)
		objectID, err := primitive.ObjectIDFromHex(electionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID"})
			return
		}
		// Query the election to ensure it is in the RESULT phase
		var election models.Election
		err = electionCollection.FindOne(context.Background(), bson.M{
			"_id":            objectID,
			"election_phase": "RESULT",
		}).Decode(&election)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Election not found in RESULT phase"})
			return
		}

		fmt.Printf("%v",election)

		// Query to find all ballots cast in the specified election
		matchBallotsQuery := bson.D{
			{"election_id", electionID},
		}

		// Count the total number of ballots in the election
		totalBallots, err := ballotCollection.CountDocuments(context.Background(), matchBallotsQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total ballots"})
			return
		}

		// Aggregate pipeline to calculate candidate-wise vote counts and percentages
		pipeline := []bson.M{
			{
				"$match": bson.M{"election_id": electionID},
			},
			{
				"$group": bson.M{
					"_id":        "$candidate_id",
					"vote_count": bson.M{"$sum": 1},
					"candidate":  bson.M{"$first": "$candidate_id"},
					"party":      bson.M{"$first": "$party_id"},
				},
			},
			{
				"$lookup": bson.M{
					"from":         "candidates",
					"localField":   "candidate",
					"foreignField": "_id",
					"as":           "candidate",
				},
			},
			{
				"$unwind": "$candidate",
			},
			{
				"$lookup": bson.M{
					"from":         "parties",
					"localField":   "party",
					"foreignField": "_id",
					"as":           "party",
				},
			},
			{
				"$unwind": "$party",
			},
			{
				"$project": bson.M{
					"_id":        0,
					"candidate":  "$candidate.username",
					"party":      "$party.name",
					"vote_count": 1,
					"percentage": bson.M{
						"$multiply": []interface{}{
							bson.M{"$divide": []interface{}{"$vote_count", totalBallots}},
							100,
						},
					},
				},
			},
			{
				"$sort": bson.M{"percentage": -1},
			},
		}

		// Aggregate results using the pipeline
		cursor, err := ballotCollection.Aggregate(context.Background(), pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate results"})
			return
		}

		// Collect and return the results
		var results []bson.M
		if err := cursor.All(context.Background(), &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve results"})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}
