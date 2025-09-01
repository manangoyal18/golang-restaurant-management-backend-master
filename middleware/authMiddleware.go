// Package middleware provides HTTP middleware functions for the restaurant management API
// This package contains authentication and authorization middleware for protecting routes
package middleware

import (
	"fmt"
	helper "golang-restaurant-management/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authentication returns a Gin middleware function that validates JWT tokens
// This middleware protects routes by requiring a valid JWT token in the request header
// The token should be provided in the "token" header field
// On success, it sets user information in the Gin context for use by handlers
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the JWT token from the "token" header
		// Client should send: headers: { "token": "jwt_token_here" }
		clientToken := c.Request.Header.Get("token")
		
		// Check if token is provided
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort() // Stop processing this request
			return
		}

		// Validate the JWT token using the helper function
		// This checks signature, expiration, and format
		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort() // Stop processing this request
			return
		}

		// Store user information from the token in the Gin context
		// This makes user data available to all subsequent handlers
		c.Set("email", claims.Email)         // User's email address
		c.Set("first_name", claims.First_name) // User's first name
		c.Set("last_name", claims.Last_name)   // User's last name
		c.Set("uid", claims.Uid)             // User's unique identifier

		// Continue to the next handler in the chain
		c.Next()
	}
}
