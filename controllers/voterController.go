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
	"github.com/mdozairq/e-voting-be/config"
	"github.com/mdozairq/e-voting-be/database"
	"github.com/mdozairq/e-voting-be/helpers"
	"github.com/mdozairq/e-voting-be/models"
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

		count, err := voterCollection.CountDocuments(ctx, bson.M{"adhaarnumber": foundAdhaar.UID})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the adhaar card"})
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

func VerifyOTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// Get the UID, phone number, and OTP from the request

		var userOtp models.OTP
		var foundVoter models.Voter

		if err := c.ShouldBindJSON(&userOtp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		utils.LogInfo(userOtp.OTP)
		utils.LogInfo(userOtp.UID)

		err := voterCollection.FindOne(ctx, bson.M{"adhaarnumber": userOtp.UID}).Decode(&foundVoter)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "adhaarcard not found, uid number seems to be incorrect"})
			return
		}

		opts := options.FindOne().SetSort(bson.M{"created_at": -1})

		// Fetch OTP document from the OTP collection based on UID and phone number
		var otpDoc models.OTP
		err = otpCollection.FindOne(ctx, bson.M{"uid": foundVoter.AdhaarNumber, "phone_number": foundVoter.Phone}, opts).Decode(&otpDoc)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "OTP not found"})
			return
		}

		// Check if OTP is expired
		if time.Now().After(otpDoc.Expiration) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired"})
			return
		}

		// Compare the user-provided OTP with the stored OTP
		if otpDoc.OTP != userOtp.OTP {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(foundVoter, "voter")

		//update tokens - token and refersh token
		helpers.UpdateAllTokens(token, refreshToken, foundVoter.ID.Hex())

		//return statusOK
		c.JSON(http.StatusOK, gin.H{"refresh_token": refreshToken, "token": token})

		// c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
	}
}

// GenerateOTP generates a 4-digit random OTP.
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

// SaveOTP saves the OTP, UID, and phone number to the OTP collection.
func SaveOTP(uid, phoneNumber, otp string) error {
	expiration := time.Now().Add(5 * time.Minute) // Set OTP expiration to 5 minutes from now
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
