package controllers

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mdozairq/e-voting-be/app/helpers"
	"github.com/mdozairq/e-voting-be/app/models"
	"github.com/mdozairq/e-voting-be/config"
	"github.com/mdozairq/e-voting-be/database"
	"github.com/mdozairq/e-voting-be/utils"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var adhaarCardCollection *mongo.Collection = database.OpenCollection(database.Client, "adhaar_cards")

// otpCollection is the MongoDB collection for storing OTP documents.
var otpCollection *mongo.Collection = database.OpenCollection(database.Client, "otps")

// voterCollection is the MongoDB collection for storing Voter Data
var voterCollection *mongo.Collection = database.OpenCollection(database.Client, "voters")

// constituencyCollection is the MongoDB collection for storing constituency
var constituencyCollection *mongo.Collection = database.OpenCollection(database.Client, "constituency")

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

func SignInVoter() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		config.LoadEnv()
		utils.LogInfo("env loaded")

		accountSid := config.NewTwilioConfig().Sid
		authToken := config.NewTwilioConfig().AuthToken

		// Initialize Twilio client
		var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: accountSid,
			Password: authToken,
		})
		var adhaar models.AdhaarCard
		var foundAdhaar models.AdhaarCard

		if err := c.BindJSON(&adhaar); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		utils.LogInfo(adhaar.UID)
		err := adhaarCardCollection.FindOne(ctx, bson.M{"uid": adhaar.UID}).Decode(&foundAdhaar)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "adhaarcard not found, uid number seems to be incorrect"})
			return
		}

		count, err := voterCollection.CountDocuments(ctx, bson.M{"adhaar_number": foundAdhaar.UID})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the adhaar card"})
			return
		}

		countStr := strconv.FormatInt(count, 10)
		utils.LogInfo(countStr)

		if count < 1 {
			CreateVoter(&foundAdhaar)
		}

		otp := GenerateOTP()

		// Save OTP to the OTP collection
		err = SaveOTP(adhaar.UID, foundAdhaar.Mobile, otp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP"})
			return
		}

		params := &api.CreateMessageParams{}
		params.SetBody("Adhaar Verification OTP:" + otp)
		params.SetFrom("+14325476702")
		params.SetTo("+91" + foundAdhaar.Mobile)

		resp, err := client.Api.CreateMessage(params)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if resp.Sid != nil {
				fmt.Println(*resp.Sid)
			} else {
				fmt.Println(resp.Sid)
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
	}

}

func AddAdhaaarCard() gin.HandlerFunc {
	return func(c *gin.Context) {

		var aadhaar models.AdhaarCard
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		if err := c.ShouldBindJSON(&aadhaar); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validate := InitializeValidator()
		if err := validate.Struct(aadhaar); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer cancel()
		count, err := adhaarCardCollection.CountDocuments(ctx, bson.M{"uid": aadhaar.UID})

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the adhaar card"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "this AadhaarCard already exsits"})
			return
		}

		fmt.Println("adhadar", aadhaar)

		dob, err := time.Parse(dobLayout, aadhaar.DOB)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date of birth format"})
			return
		}

		aadhaar.DOB = dob.Format(dobLayout)

		_, insertErr := adhaarCardCollection.InsertOne(ctx, aadhaar)
		if insertErr != nil {
			msg := fmt.Sprintf("nexVoter item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		createConstituency(&aadhaar)

		c.JSON(http.StatusCreated, gin.H{"message": "Aadhaar card created successfully"})

	}
}

func createConstituency(aadhaarCard *models.AdhaarCard) {
	count, err := constituencyCollection.CountDocuments(context.Background(), bson.M{"district": aadhaarCard.City, "state": aadhaarCard.State, "country": aadhaarCard.Country})

	if err != nil {
		log.Panic(err)
		return
	}

	if count > 0 {
		return
	}

	constituencyDoc := models.Constituency{
		ID:       primitive.NewObjectID(),
		City:     aadhaarCard.City,
		District: aadhaarCard.City,
		State:    aadhaarCard.State,
		Country:  aadhaarCard.Country,
		PinCode:  aadhaarCard.PinCode,
	}

	_, err = constituencyCollection.InsertOne(context.Background(), constituencyDoc)
	if err != nil {
		log.Panic(err)
		return
	}
	return
}

func VerifyOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var userOtp models.OTP
		var foundVoter models.Voter

		if err := c.ShouldBindJSON(&userOtp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		utils.LogInfo(userOtp.OTP)
		utils.LogInfo(userOtp.UID)

		err := voterCollection.FindOne(ctx, bson.M{"adhaar_number": userOtp.UID}).Decode(&foundVoter)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "adhaarcard not found, uid number seems to be incorrect"})
			return
		}

		opts := options.FindOne().SetSort(bson.M{"created_at": -1})

		var otpDoc models.OTP
		err = otpCollection.FindOne(ctx, bson.M{"uid": foundVoter.AdhaarNumber, "phone_number": foundVoter.Phone}, opts).Decode(&otpDoc)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "OTP not found"})
			return
		}

		if time.Now().After(otpDoc.Expiration) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired"})
			return
		}

		if otpDoc.OTP != userOtp.OTP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
			return
		}

		token, refreshToken, _ := helpers.GenerateVoterTokens(foundVoter, "voter")

		helpers.UpdateAllTokens(token, refreshToken, foundVoter.ID.Hex())

		c.JSON(http.StatusOK, gin.H{"refresh_token": refreshToken, "token": token})
	}
}

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func SaveOTP(uid, phoneNumber, otp string) error {
	expiration := time.Now().Add(5 * time.Minute)
	createdAt := time.Now()

	otpDoc := models.OTP{
		UID:         uid,
		PhoneNumber: phoneNumber,
		OTP:         otp,
		Expiration:  expiration,
		CreatedAt:   createdAt,
	}

	_, err := otpCollection.InsertOne(context.Background(), otpDoc)
	return err
}

func CreateVoter(aadhaarCard *models.AdhaarCard) error {
	createdAt := time.Now()
	updatedAt := time.Now()
	const votingAge = 18
	dob, err := time.Parse("2006-01-02", aadhaarCard.DOB)
	if err != nil {
		fmt.Println("Error parsing DOB:", err)
	}
	// Calculate age
	age := calculateAge(dob)

	eligibleForVoting := age >= votingAge

	voterDoc := models.Voter{
		ID:           primitive.NewObjectID(),
		Name:         aadhaarCard.Name,
		AdhaarNumber: aadhaarCard.UID,
		Phone:        aadhaarCard.Mobile,
		DateOfBirth:  aadhaarCard.DOB,
		Gender:       aadhaarCard.Gender,
		IsEligible:   eligibleForVoting,
		IsVoted:      false,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	_, err = voterCollection.InsertOne(context.Background(), voterDoc)
	return err
}

func calculateAge(dob time.Time) int {
	today := time.Now()
	years := today.Year() - dob.Year()
	if today.YearDay() < dob.YearDay() {
		years--
	}
	return years
}

func GetElectionByAadhaarLocation() gin.HandlerFunc {
	return func(c *gin.Context) {
		adhaarNumber := c.Query("adhaarNumber")

		// Query Aadhaar card data based on the adhaar number
		var aadhaar models.AdhaarCard
		err := adhaarCardCollection.FindOne(context.Background(), bson.M{"uid": adhaarNumber}).Decode(&aadhaar)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Aadhaar card data"})
			return
		}

		// Query Constituency data based on the Aadhaar card state, city, and country
		var constituency models.Constituency
		err = constituencyCollection.FindOne(context.Background(), bson.M{
			"state":   aadhaar.State,
			"city":    aadhaar.City,
			"country": aadhaar.Country,
		}).Decode(&constituency)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch Constituency data"})
			return
		}
		log.Printf("Constituency: %+v", constituency)
		// Query Election based on the fetched Constituency data and ElectionPhase
		var election models.Election
		err = electionCollection.FindOne(context.Background(), bson.M{
			"constituency":   constituency.ID.Hex(),
			"election_phase": "VOTING",
		}).Decode(&election)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No Election found for the specified location and phase"})
			return
		}
		log.Printf("Election: %+v", election)
		// Query candidates registered for the specified election
		var registeredCandidates []models.Candidate
		candidateCursor, err := candidateCollection.Find(context.Background(), bson.D{
			{"election_id", election.ID.Hex()},
			{"is_registered", true},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch registered candidates"})
			return
		}
		defer candidateCursor.Close(context.Background())
		for candidateCursor.Next(context.Background()) {
			var candidate models.Candidate
			if err := candidateCursor.Decode(&candidate); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode candidate"})
				return
			}

			// Populate Party data for each registered candidate
			var party models.Party
			err := partyCollection.FindOne(context.Background(), bson.M{"_id": candidate.PartyID}).Decode(&party)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch party data for candidate"})
				return
			}
			log.Printf("Party: %+v", party)
			candidate.Party = party
			registeredCandidates = append(registeredCandidates, candidate)
		}

		c.JSON(http.StatusOK, gin.H{
			"election":             election,
			"registeredCandidates": registeredCandidates,
		})
	}
}
