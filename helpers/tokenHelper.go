// Package helper provides utility functions for JWT token management
// This package handles token generation, validation, and database updates for user authentication
package helper

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SignedDetails represents the JWT token payload structure
// This struct contains user information and standard JWT claims
// It embeds jwt.StandardClaims for expiration and other standard fields
type SignedDetails struct {
	// Email is the user's email address stored in the token
	Email string
	// First_name is the user's first name stored in the token
	First_name string
	// Last_name is the user's last name stored in the token
	Last_name string
	// Uid is the user's unique identifier stored in the token
	Uid string
	// StandardClaims provides standard JWT fields like expiration time
	jwt.StandardClaims
}

// userCollection is a global reference to the MongoDB "user" collection
// Used for updating user tokens in the database
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// SECRET_KEY is the JWT signing key retrieved from environment variables
// This key is used to sign and validate all JWT tokens
var SECRET_KEY string = os.Getenv("SECRET_KEY")

// GenerateAllTokens creates both access and refresh JWT tokens for a user
// Parameters:
//   - email: user's email address
//   - firstName: user's first name
//   - lastName: user's last name  
//   - uid: user's unique identifier
// Returns: access token, refresh token, and any error
func GenerateAllTokens(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {
	// Create claims for the access token (expires in 24 hours)
	// Contains user information for API authorization
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			// Access token expires in 24 hours
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	// Create claims for the refresh token (expires in 7 days)
	// Contains minimal information, used only for token renewal
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			// Refresh token expires in 168 hours (7 days)
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	// Generate the signed access token using HS256 algorithm
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	// Generate the signed refresh token using HS256 algorithm
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

// UpdateAllTokens updates both access and refresh tokens for a user in the database
// This function is called after login or token refresh to store new tokens
// Parameters:
//   - signedToken: the new access token to store
//   - signedRefreshToken: the new refresh token to store
//   - userId: the user's unique identifier
func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	// Set up context with timeout for the database operation
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	// Build the update document with new token values
	var updateObj primitive.D

	// Add the new access token to the update object
	updateObj = append(updateObj, bson.E{"token", signedToken})
	// Add the new refresh token to the update object
	updateObj = append(updateObj, bson.E{"refresh_token", signedRefreshToken})

	// Update the "updated_at" timestamp to track when tokens were refreshed
	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at", Updated_at})

	// Configure upsert option - creates document if it doesn't exist
	upsert := true
	// Filter to find the user document by user_id
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	// Perform the database update operation
	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj}, // Use $set operator to update specific fields
		},
		&opt,
	)
	// Clean up the context resources
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}
	return
}

// ValidateToken parses and validates a JWT token
// This function checks token signature, format, and expiration
// Parameters:
//   - signedToken: the JWT token string to validate
// Returns:
//   - claims: the parsed token claims if valid
//   - msg: error message if validation fails, empty if success
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	// Parse the token with custom claims structure
	// The keyFunc returns the secret key used to verify the token signature
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{}, // Expected claims structure
		func(token *jwt.Token) (interface{}, error) {
			// Return the secret key for signature validation
			return []byte(SECRET_KEY), nil
		},
	)

	// Extract claims from the parsed token
	// Check if the token claims can be cast to our custom SignedDetails type
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		// Token format is invalid or claims structure doesn't match
		msg = fmt.Sprintf("the token is invalid")
		if err != nil {
			msg = err.Error()
		}
		return
	}

	// Check if the token has expired
	// Compare the token's expiration time with current time
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		if err != nil {
			msg = err.Error()
		}
		return
	}

	// Token is valid - return claims with no error message
	return claims, msg
}
