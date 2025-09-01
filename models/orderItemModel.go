package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderItem represents an individual item within a customer order
// This struct defines the structure of order item documents stored in MongoDB
// Each order item links a food item to an order with quantity and pricing info
type OrderItem struct {
	// ID is the MongoDB ObjectID - the unique identifier for the order item document
	ID primitive.ObjectID `bson:"_id"`
	
	// Quantity represents the size/quantity of the food item (required: S, M, or L)
	// Note: This appears to be for portion sizes rather than numeric quantity
	Quantity *string `json:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	
	// Unit_price is the price for this specific order item (required)
	// This may differ from the base food price due to size or modifications
	Unit_price *float64 `json:"unit_price" validate:"required"`
	
	// Created_at is the timestamp when the order item was added to the order
	Created_at time.Time `json:"created_at"`
	
	// Updated_at is the timestamp when the order item was last modified
	Updated_at time.Time `json:"updated_at"`
	
	// Food_id is the reference to the food item being ordered (required)
	// This creates a relationship between order items and food items
	Food_id *string `json:"food_id" validate:"required"`
	
	// Order_item_id is the string representation of the MongoDB ObjectID
	// Used for easier referencing in other collections and API responses
	Order_item_id string `json:"order_item_id"`
	
	// Order_id is the reference to the parent order (required)
	// This creates a relationship between order items and their parent order
	Order_id string `json:"order_id" validate:"required"`
}
