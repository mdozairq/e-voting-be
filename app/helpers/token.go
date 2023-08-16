package helpers

import (
	// "context"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mdozairq/e-voting-be/app/models"
	"github.com/mdozairq/e-voting-be/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type VoterSignedDetails struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	Name         string             `bson:"name" json:"name" validate:"required,min=2,max=100"`
	AdhaarNumber string             `bson:"adhaar_number" json:"adhaar_number" validate:"required,min=12"`
	Phone        string             `bson:"phone" json:"phone" validate:"regexp=(0|+91|091|91)[0-9]+$"`
	DateOfBirth  string             `bson:"date_of_birth" json:"date_of_birth" validate:"required"`
	Gender       string             `bson:"gender" json:"gender" validate:"required,eq=MALE|eq=FEMALE"`
	IsEligible   bool               `bson:"is_eligible" json:"is_eligible" validate:"required"`
	IsVoted      bool               `bson:"is_voted" json:"is_voted" validate:"required"`
	Role         string             `json:"role"`
	jwt.StandardClaims
}

type AdminDetails struct {
	AdminToken string
	Email      string
	Role       string
	jwt.StandardClaims
}

type CandidateClaims struct {
	CandidateID    primitive.ObjectID `json:"candidate_id"`
	VoterID        primitive.ObjectID `json:"voter_id"`
	Name           string             `json:"name"`
	AdhaarNumber   string             `json:"adhaar_number"`
	Phone          string             `json:"phone"`
	Gender         string             `json:"gender"`
	IsEligible     bool               `json:"is_eligible"`
	IsVoted        bool               `json:"is_voted"`
	IsRegistered   bool               `json:"is_registered"`
	ElectionID     string             `json:"election_id"`
	PartyID        string             `json:"party_id"`
	HasCrimeRecord bool               `json:"has_crime_record"`
	Role           string             `json:"role"`
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAdminToken(email string, adminAuthToken string, role string) (signedToken string, err error) {
	AdminClaim := &AdminDetails{
		AdminToken: adminAuthToken,
		Email:      email,
		Role:       role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AdminClaim)

	// Sign the token with the admin's authentication token (secret)
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func GenerateVoterTokens(voter models.Voter, role string) (signedToken string, signedRefreshToken string, err error) {
	claims := &VoterSignedDetails{
		ID:           voter.ID,
		Name:         voter.Name,
		AdhaarNumber: voter.AdhaarNumber,
		Phone:        voter.Phone,
		DateOfBirth:  voter.DateOfBirth,
		Gender:       voter.Gender,
		IsEligible:   voter.IsEligible,
		IsVoted:      voter.IsVoted,
		Role:         role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &VoterSignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err

}

func GenerateCandidateTokens(candidate models.Candidate, voter models.Voter, role string) (signedToken string, signedRefreshToken string, err error) {
	claims := &CandidateClaims{
		CandidateID:    candidate.ID,
		VoterID:        voter.ID,
		Name:           voter.Name,
		AdhaarNumber:   voter.AdhaarNumber,
		Phone:          voter.Phone,
		Gender:         voter.Gender,
		IsEligible:     voter.IsEligible,
		IsVoted:        voter.IsVoted,
		IsRegistered:   candidate.IsRegistered,
		ElectionID:     candidate.ElectionID,
		PartyID:        candidate.PartyID,
		HasCrimeRecord: candidate.HasCrimeRecords,
		Role:           role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &CandidateClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err

}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token", signedToken})
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	return
}

func ValidateToken(signedToken string) (claims *VoterSignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&VoterSignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	claims, ok := token.Claims.(*VoterSignedDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	//the token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg

}

func ValidateCandidateToken(signedToken string) (claims *CandidateClaims, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&CandidateClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	claims, ok := token.Claims.(*CandidateClaims)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	//the token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg

}

func AdminValidateToken(signedToken string) (claims *AdminDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&AdminDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	claims, ok := token.Claims.(*AdminDetails)
	if !ok {
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	//the token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		msg = err.Error()
		return
	}

	return claims, msg

}
