package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mdozairq/e-voting-be/database"
	"github.com/mdozairq/e-voting-be/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
