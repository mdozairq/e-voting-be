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
	"github.com/mdozairq/e-voting-be/app/helpers"
	"github.com/mdozairq/e-voting-be/app/models"
	"github.com/mdozairq/e-voting-be/database"
	"github.com/mdozairq/e-voting-be/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var candidateCollection *mongo.Collection = database.OpenCollection(database.Client, "candidate")
var validate = validator.New()

type CandidateDto struct {
	Username     string `json:"username" validate:"required"`
	AdhaarNumber string `json:"uid" validate:"required,len=12"`
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=8"`
}

type SignInCandidateDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

func GetCandidates() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		pageSize, err := strconv.Atoi(c.Query("page_size"))
		if err != nil || pageSize < 1 {
			pageSize = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * pageSize
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"candidates", bson.D{{"$slice", []interface{}{"$data", startIndex, pageSize}}}},
				}}}

		result, err := candidateCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while listing candidates"})
			return
		}

		defer cancel()

		var allCandidates []bson.M
		if err = result.All(ctx, &allCandidates); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allCandidates[0])
	}
}

func SignUpCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var requestedCandidate CandidateDto
		var voter models.Voter

		//convert the JSON data coming from postman to something that golang understands
		if err := c.BindJSON(&requestedCandidate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//validate the data based on candidate struct
		validationErr := validate.Struct(&requestedCandidate)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := voterCollection.FindOne(ctx, bson.M{"adhaar_number": requestedCandidate.AdhaarNumber}).Decode(&voter)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "voter_id not found, voter not seems to be registered"})
			return
		}

		if !voter.IsEligible {
			c.JSON(http.StatusBadRequest, gin.H{"error": "voter is not eligible or under age"})
			return
		}

		count, err := candidateCollection.CountDocuments(ctx, bson.M{"voter_id": voter.ID})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the voterId"})
			return
		}

		//hash password
		password := HashPassword(requestedCandidate.Password)

		//you'll also check if the phone no. has already been used by another candidate
		count, err = candidateCollection.CountDocuments(ctx, bson.M{"username": requestedCandidate.Username})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for the username"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this Candidate with same adhaar_number or username already exists"})
			return
		}

		createdAt := time.Now()
		updatedAt := time.Now()
		candidateID := primitive.NewObjectID()

		createCandidate := models.Candidate{
			ID:                       candidateID,
			Email:                    requestedCandidate.Email,
			Username:                 requestedCandidate.Username,
			VoterID:                  voter.ID.Hex(),
			Phone:                    voter.Phone,
			Password:                 password,
			PartyID:                  "",
			RegisteredConstituencyID: "",
			Assets:                   "",
			HasCrimeRecords:          false,
			IsAccused:                false,
			IsEligible:               voter.IsEligible,
			IsRegistered:             false,
			CreatedAt:                createdAt,
			UpdatedAt:                updatedAt,
		}

		//if all ok, then you insert this new candidate into the candidate collection
		_, insertErr := candidateCollection.InsertOne(ctx, createCandidate)
		if insertErr != nil {
			msg := fmt.Sprintf("candidate item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		token, refreshToken, _ := helpers.GenerateCandidateTokens(createCandidate, voter, "candidate")

		helpers.UpdateAllTokens(token, refreshToken, voter.ID.Hex())

		c.JSON(http.StatusOK, gin.H{"refresh_token": refreshToken, "token": token})
	}
}

func SignInCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var requestedCandidate SignInCandidateDto
		var foundCandidate models.Candidate
		var voter models.Voter

		if err := c.BindJSON(&requestedCandidate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(&requestedCandidate)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		err := candidateCollection.FindOne(ctx, bson.M{"username": requestedCandidate.Username}).Decode(&foundCandidate)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found, user seems to be incorrect"})
			return
		}

		isValidPassword, msg := VerifyPassword(foundCandidate.Password, requestedCandidate.Password)
		if !isValidPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}
		utils.LogInfo(foundCandidate.VoterID)
		voterId, err := primitive.ObjectIDFromHex(foundCandidate.VoterID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert VoterID to ObjectID"})
			return
		}
		err = voterCollection.FindOne(ctx, bson.M{"_id": voterId}).Decode(&voter)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Voter not found, voter_id number seems to be incorrect"})
			return
		}

		token, refreshToken, _ := helpers.GenerateCandidateTokens(foundCandidate, voter, "candidate")

		helpers.UpdateAllTokens(token, refreshToken, voter.ID.Hex())

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func GetCandidateByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the candidate ID from the request parameters
		candidateID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(candidateID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
			return
		}

		// Create a MongoDB context
		ctx := context.Background()

		// Query the candidate by their ID
		var candidate models.Candidate
		err = candidateCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&candidate)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found"})
			return
		}

		c.JSON(http.StatusOK, candidate)
	}
}

func UpdateCandidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the candidate ID from the request parameters
		candidateID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(candidateID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
			return
		}

		// Parse the updated candidate data from the request body
		var updatedCandidate models.Candidate
		if err := c.ShouldBindJSON(&updatedCandidate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set the updated timestamp
		updatedCandidate.UpdatedAt = time.Now()

		// Update the candidate in the MongoDB collection
		ctx := context.Background()
		update := bson.M{"$set": updatedCandidate}
		_, err = candidateCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
		if err != nil {
			print()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update candidate"})
			return
		}

		c.JSON(http.StatusOK, updatedCandidate)
	}
}

func UpdateCandidatePartial() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the candidate ID from the request parameters
		candidateID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(candidateID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
			return
		}

		// Parse the partial candidate data from the request body
		var partialUpdate models.Candidate
		if err := c.ShouldBindJSON(&partialUpdate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create an update with only the fields you want to update
		update := bson.M{}
		if partialUpdate.Assets != "" {
			update["assets"] = partialUpdate.Assets
		}
		if partialUpdate.HasCrimeRecords {
			update["has_crime_records"] = partialUpdate.HasCrimeRecords
		}
		if partialUpdate.IsAccused {
			update["is_accused"] = partialUpdate.IsAccused
		}
		if partialUpdate.RegisteredConstituencyID != "" {
			update["registered_constituency_id"] = partialUpdate.RegisteredConstituencyID
		}
		if partialUpdate.ElectionID != "" {
			update["election_id"] = partialUpdate.ElectionID
		}
		if partialUpdate.PartyID != "" {
			update["party_id"] = partialUpdate.PartyID
		}
		update["updated_at"] = time.Now()
		update["is_registered"] = true
		// Update the candidate in the MongoDB collection
		ctx := context.Background()
		_, err = candidateCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update candidate"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Candidate updated successfully"})
	}
}

func GetCandidatesByElectionID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the election ID from the query parameter
		electionID := c.Query("electionId")

		// Query to find candidates associated with the specified election
		matchCandidatesQuery := bson.D{
			{"election_id", electionID},
		}

		// Find candidate documents based on the query
		var candidates []models.Candidate
		cursor, err := candidateCollection.Find(context.Background(), matchCandidatesQuery)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch candidates"})
			return
		}
		defer cursor.Close(context.Background())
		// Collect candidates that match the query
		for cursor.Next(context.Background()) {
			var candidate models.Candidate
			if err := cursor.Decode(&candidate); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode candidate"})
				return
			}
			log.Printf("Candy: %+v", candidate)
			// Populate Party data for each registered candidate
			objectID, err := primitive.ObjectIDFromHex(candidate.PartyID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid constituency ID"})
				return
			}

			var party models.Party
			err = partyCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&party)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch party data for candidate"})
				return
			}
			log.Printf("Party: %+v", party)
			candidate.Party = party

			// Populate Voter data for each registered candidate
			objectID, err = primitive.ObjectIDFromHex(candidate.VoterID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid voter ID"})
				return
			}

			var voter models.Voter
			err = voterCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&voter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch voter data for candidate"})
				return
			}

			candidate.Voter = voter

			candidates = append(candidates, candidate)
		}

		c.JSON(http.StatusOK, gin.H{"candidates": candidates})
	}
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	utils.LogInfo(userPassword)

	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}
	return check, msg
}
