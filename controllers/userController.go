// Package controller handles all HTTP request processing for user management
// This package contains handlers for user registration, authentication, and user data operations
package controller

import (
	"context"
	"fmt"
	"golang-restaurant-management/database"
	helper "golang-restaurant-management/helpers"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// userCollection is a global reference to the MongoDB "user" collection
// This provides access to user documents in the database
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// validate is a global validator instance for struct validation
// This is used to validate incoming JSON data against struct validation tags
var validate = validator.New()

// GetUsers returns a gin handler function that retrieves a paginated list of users
// This endpoint supports pagination with query parameters: page, recordPerPage, startIndex
// Returns: JSON array of user objects with pagination metadata
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set up context with timeout to prevent long-running database queries
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Parse pagination parameters from query string
		// recordPerPage determines how many users to return per page (default: 10)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		// page determines which page of results to return (default: 1)
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		// Calculate the starting index for pagination
		startIndex := (page - 1) * recordPerPage
		// Allow override via query parameter if provided
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		// MongoDB aggregation pipeline stages
		// matchStage: matches all documents (empty filter)
		matchStage := bson.D{{"$match", bson.D{{}}}}
		// projectStage: shapes the output and implements pagination using $slice
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}

		// Execute the aggregation pipeline to get paginated results
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}

		// Parse the aggregation result into a slice of documents
		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		// Return the first (and only) result from the aggregation
		c.JSON(http.StatusOK, allUsers[0])

	}
}

// GetUser returns a gin handler function that retrieves a specific user by ID
// Parameters: user_id (URL parameter) - the unique identifier of the user to retrieve
// Returns: JSON object containing the user's information
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set up context with timeout for the database operation
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// Extract user_id from the URL parameters
		userId := c.Param("user_id")

		// Create an empty User struct to hold the database result
		var user models.User

		// Find the user document by user_id field and decode it into the User struct
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		// Clean up the context resources
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		// Return the user data as JSON
		c.JSON(http.StatusOK, user)
	}
}

// SignUp returns a gin handler function for user registration
// This endpoint creates a new user account with email/phone uniqueness validation
// Request body should contain user details in JSON format
// Returns: JSON object with insertion result or error message
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set up context with timeout for database operations
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		// Parse JSON request body into User struct
		// This converts the JSON data coming from client to Go struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the user data against struct validation tags
		// This checks required fields, email format, string lengths, etc.
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Check if the email has already been used by another user
		// Ensures email uniqueness across all user accounts
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		// Hash the password using bcrypt before storing it
		// This ensures passwords are never stored in plain text
		password := HashPassword(*user.Password)
		user.Password = &password

		// Check if the phone number has already been used by another user
		// Ensures phone uniqueness across all user accounts
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}

		// If email or phone already exists, return error
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exsits"})
			return
		}

		// Set up timestamps and ID for the new user document
		// Parse current time into RFC3339 format for consistency
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		// Generate new MongoDB ObjectID and convert to string
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// Generate JWT access and refresh tokens for the new user
		// This allows immediate login after registration
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		// Insert the new user document into the MongoDB collection
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		// Return success response with the insertion result
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

// Login returns a gin handler function for user authentication
// This endpoint validates user credentials and returns JWT tokens
// Request body should contain email and password in JSON format
// Returns: JSON object with user data and updated tokens
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set up context with timeout for database operations
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User      // Holds login request data
		var foundUser models.User // Holds user data from database

		// Parse JSON login request body into User struct
		// This converts the login data from client to Go struct
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find user by email in the database
		// This checks if a user with the provided email exists
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found, login seems to be incorrect"})
			return
		}

		// Verify the provided password against the hashed password in database
		// This uses bcrypt to compare the plain text password with the hash
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// Generate new JWT access and refresh tokens for the authenticated user
		// This creates fresh tokens for the session
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		// Update the user's tokens in the database
		// This ensures the latest tokens are stored for future validation
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		// Return success response with user data (tokens are updated in database)
		c.JSON(http.StatusOK, foundUser)
	}
}

// HashPassword takes a plain text password and returns a bcrypt hash
// Parameters: password (string) - the plain text password to hash
// Returns: string - the bcrypt hashed password
// The cost factor of 14 provides good security while maintaining reasonable performance
func HashPassword(password string) string {
	// Generate bcrypt hash with cost factor of 14
	// Higher cost means more secure but slower hashing
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	// Convert byte array to string and return
	return string(bytes)
}

// VerifyPassword compares a plain text password with a bcrypt hash
// Parameters:
//   - userPassword: the plain text password provided by user
//   - providedPassword: the hashed password stored in database
// Returns:
//   - bool: true if passwords match, false otherwise
//   - string: error message if verification fails, empty string if success
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	// Use bcrypt to compare the plain text password with the hash
	// bcrypt handles the salt and hashing internally
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	// If comparison fails, set error message and return false
	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}
	// Return verification result and any error message
	return check, msg
}
