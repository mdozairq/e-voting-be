package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mdozairq/e-voting-be/database"
	"github.com/mdozairq/e-voting-be/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var adhaarCardCollection *mongo.Collection = database.OpenCollection(database.Client, "adhaarcards")

const dobLayout = "2006-01-02"

// customDateValidator validates that the provided string is a valid date in the format "YYYY-MM-DD".
func customDateValidator(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse(dobLayout, dateStr)
	return err == nil
}

// InitializeValidator initializes the custom validators for the validator.
func InitializeValidator() *validator.Validate {
	v := validator.New()
	if err := v.RegisterValidation("date", customDateValidator); err != nil {
		panic("error registering validation function")
	}
	return v
}

func GetVoters() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetVoter() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SignUpVoter() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SignInVoter() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func AddAdhaaarCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var existingAdhaar models.AdhaarCard
		var aadhaar models.AdhaarCard
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		if err := c.ShouldBindJSON(&aadhaar); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate Aadhaar card data
		validate := InitializeValidator()
		if err := validate.Struct(aadhaar); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		count, err := adhaarCardCollection.CountDocuments(ctx, bson.M{"uid": aadhaar.UID})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "this AadhaarCard already exsits"})
			return
		}

		fmt.Println("adhadar", aadhaar)
		// Parse DOB from the custom layout
		dob, err := time.Parse(dobLayout, aadhaar.DOB)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date of birth format"})
			return
		}

		aadhaar.DOB = dob.Format(dobLayout)

		// Insert Aadhaar card data into MongoDB
		_, insertErr := adhaarCardCollection.InsertOne(ctx, aadhaar)
		if insertErr != nil {
			msg := fmt.Sprintf("nexVoter item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusCreated, gin.H{"message": "Aadhaar card created successfully"})

	}
}
