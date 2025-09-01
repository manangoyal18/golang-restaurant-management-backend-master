// Package routes defines all HTTP route handlers for the restaurant management API
// This package organizes routes by functionality (users, food, orders, etc.)
package routes

import (
	controller "golang-restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

// UserRoutes sets up all user-related HTTP routes
// These routes handle user authentication, registration, and user data management
// Note: signup and login routes are public (no authentication required)
// Other user routes require authentication middleware to be applied externally
func UserRoutes(incomingRoutes *gin.Engine) {
	// GET /users - Retrieve paginated list of all users
	// Requires authentication (applied in main.go)
	incomingRoutes.GET("/users", controller.GetUsers())
	
	// GET /users/:user_id - Retrieve specific user by ID
	// Requires authentication (applied in main.go)
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	
	// POST /users/signup - Register a new user account
	// Public route - no authentication required
	incomingRoutes.POST("/users/signup", controller.SignUp())
	
	// POST /users/login - Authenticate user and receive JWT tokens
	// Public route - no authentication required
	incomingRoutes.POST("/users/login", controller.Login())
}
