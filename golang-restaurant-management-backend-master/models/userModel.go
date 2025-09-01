// Package models defines the data structures used throughout the restaurant management system
// This package contains all the MongoDB document models and their validation rules
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user account in the restaurant management system
// This struct defines the structure of user documents stored in MongoDB
// It includes validation tags for input validation and JSON tags for API responses
type User struct {
	// ID is the MongoDB ObjectID - the unique identifier for the user document
	ID primitive.ObjectID `bson:"_id"`
	
	// First_name is the user's first name (required, 2-100 characters)
	// Using pointer to string allows for null values and better memory management
	First_name *string `json:"first_name" validate:"required,min=2,max=100"`
	
	// Last_name is the user's last name (required, 2-100 characters)
	Last_name *string `json:"last_name" validate:"required,min=2,max=100"`
	
	// Password is the user's hashed password (required, minimum 6 characters before hashing)
	// This will be hashed using bcrypt before storing in the database
	Password *string `json:"Password" validate:"required,min=6"`
	
	// Email is the user's email address (required, must be valid email format)
	// This serves as a unique identifier for login purposes
	Email *string `json:"email" validate:"email,required"`
	
	// Avatar is an optional URL to the user's profile picture
	Avatar *string `json:"avatar"`
	
	// Phone is the user's phone number (required)
	// This can be used for notifications and account verification
	Phone *string `json:"phone" validate:"required"`
	
	// Token is the JWT access token for authentication
	// This is generated when the user logs in and used for API requests
	Token *string `json:"token"`
	
	// Refresh_Token is the JWT refresh token for token renewal
	// This allows generating new access tokens without requiring re-login
	Refresh_Token *string `json:"refresh_token"`
	
	// Created_at is the timestamp when the user account was created
	Created_at time.Time `json:"created_at"`
	
	// Updated_at is the timestamp when the user account was last modified
	Updated_at time.Time `json:"updated_at"`
	
	// User_id is the string representation of the MongoDB ObjectID
	// This is used for easier referencing in other collections and API responses
	User_id string `json:"user_id"`
}
