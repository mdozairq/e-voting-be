package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mdozairq/e-voting-be/database"
	// "github.com/mdozairq/e-voting-be/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var candidateCollection *mongo.Collection = database.OpenCollection(database.Client, "candidate")
var validate = validator.New()

func Getcandidates() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		page_size, err := strconv.Atoi(c.Query("page_size"))
		if err != nil || page_size < 1 {
			page_size = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * page_size
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, page_size}}}},
				}}}

		result, err := candidateCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing food items"})
		}
		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allFoods[0])
	}
}

func SignUpCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SignInCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func HashPassword(password string) string {
	return password
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	return userPassword == providedPassword,
		userPassword
}
