package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Food represents a food item in the restaurant management system
// This struct defines the structure of food documents stored in MongoDB
// Each food item belongs to a menu and has pricing and image information
type Food struct {
	// ID is the MongoDB ObjectID - the unique identifier for the food document
	ID primitive.ObjectID `bson:"_id"`
	
	// Name is the food item's name (required, 2-100 characters)
	// This is the display name shown to customers
	Name *string `json:"name" validate:"required,min=2,max=100"`
	
	// Price is the cost of the food item (required)
	// Stored as float64 to handle decimal currency values
	Price *float64 `json:"price" validate:"required"`
	
	// Food_image is the URL or path to the food item's image (required)
	// Used for displaying the food item visually in menus and orders
	Food_image *string `json:"food_image" validate:"required"`
	
	// Created_at is the timestamp when the food item was added
	Created_at time.Time `json:"created_at"`
	
	// Updated_at is the timestamp when the food item was last modified
	Updated_at time.Time `json:"updated_at"`
	
	// Food_id is the string representation of the MongoDB ObjectID
	// Used for easier referencing in other collections and API responses
	Food_id string `json:"food_id"`
	
	// Menu_id is the reference to the menu this food item belongs to (required)
	// This creates a relationship between food items and their parent menu
	Menu_id *string `json:"menu_id" validate:"required"`
}
