package helpers

import (
	// "context"
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mdozairq/e-voting-be/database"
	"github.com/mdozairq/e-voting-be/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type VoterSignedDetails struct {
	ID           primitive.ObjectID
	Name         string
	AdhaarNumber string
	Phone        string
	DateOfBirth  string
	Gender       string
	IsEligible   bool
	IsVoted      bool
	Role         string
	jwt.StandardClaims
}

type AdminDetails struct {
	AdminToken string
	Email      string
	Role       string
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

func GenerateAllTokens(voter models.Voter, role string) (signedToken string, signedRefreshToken string, err error) {
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