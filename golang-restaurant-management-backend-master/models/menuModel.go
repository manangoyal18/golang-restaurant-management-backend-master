package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Menu represents a restaurant menu in the management system
// This struct defines the structure of menu documents stored in MongoDB
// Menus can have time-based availability and categorization
type Menu struct {
	// ID is the MongoDB ObjectID - the unique identifier for the menu document
	ID primitive.ObjectID `bson:"_id"`
	
	// Name is the menu's display name (required)
	// Examples: "Breakfast Menu", "Dinner Specials", "Weekend Brunch"
	Name string `json:"name" validate:"required"`
	
	// Category is the menu classification (required)
	// Examples: "breakfast", "lunch", "dinner", "desserts", "beverages"
	Category string `json:"category" validate:"required"`
	
	// Start_Date is when this menu becomes available (optional)
	// Used for seasonal or time-limited menus
	Start_Date *time.Time `json:"start_date"`
	
	// End_Date is when this menu expires (optional)
	// Used for seasonal or promotional menus
	End_Date *time.Time `json:"end_date"`
	
	// Created_at is the timestamp when the menu was created
	Created_at time.Time `json:"created_at"`
	
	// Updated_at is the timestamp when the menu was last modified
	Updated_at time.Time `json:"updated_at"`
	
	// Menu_id is the string representation of the MongoDB ObjectID
	// Used for easier referencing in other collections and API responses
	Menu_id string `json:"food_id"`
}
