// Package main is the entry point for the restaurant management backend API server
// This application provides a complete REST API for managing restaurant operations
// including users, food items, menus, tables, orders, and invoices
package main

import (
	"os"

	"golang-restaurant-management/database"

	middleware "golang-restaurant-management/middleware"
	routes "golang-restaurant-management/routes"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/mongo"
)

// foodCollection is a global MongoDB collection reference for food items
// This provides access to the "food" collection in the "restaurant" database
var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

// main is the entry point of the application
// It sets up the HTTP server with middleware, routes, and starts listening on the specified port
func main() {
	// Get the port number from environment variable, default to 8000 if not set
	// This allows for flexible deployment configurations
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	// Create a new Gin router instance
	// Gin is a HTTP web framework for Go that provides fast routing and middleware support
	router := gin.New()
	
	// Add logging middleware to log HTTP requests
	// This helps with debugging and monitoring API usage
	router.Use(gin.Logger())
	
	// Set up user routes (login, signup) - these don't require authentication
	// User routes are public endpoints for registration and authentication
	routes.UserRoutes(router)
	
	// Apply authentication middleware to all subsequent routes
	// This ensures that all routes below this line require a valid JWT token
	router.Use(middleware.Authentication())

	// Set up protected routes that require authentication
	// These routes handle the core restaurant management functionality
	routes.FoodRoutes(router)        // CRUD operations for food items
	routes.MenuRoutes(router)        // Menu management endpoints
	routes.TableRoutes(router)       // Table management for restaurant seating
	routes.OrderRoutes(router)       // Order processing and management
	routes.OrderItemRoutes(router)   // Individual order item management
	routes.InvoiceRoutes(router)     // Invoice generation and management

	// Start the HTTP server on the specified port
	// The server will listen for incoming HTTP requests and route them appropriately
	router.Run(":" + port)
}
