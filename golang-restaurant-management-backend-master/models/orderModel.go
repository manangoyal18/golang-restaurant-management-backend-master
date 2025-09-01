package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order represents a customer order in the restaurant management system
// This struct defines the structure of order documents stored in MongoDB
// Each order is associated with a table and contains timing information
type Order struct {
	// ID is the MongoDB ObjectID - the unique identifier for the order document
	ID primitive.ObjectID `bson:"_id"`
	
	// Order_Date is when the order was placed (required)
	// This timestamp is used for order tracking and reporting
	Order_Date time.Time `json:"order_date" validate:"required"`
	
	// Created_at is the timestamp when the order record was created in the system
	Created_at time.Time `json:"created_at"`
	
	// Updated_at is the timestamp when the order was last modified
	Updated_at time.Time `json:"updated_at"`
	
	// Order_id is the string representation of the MongoDB ObjectID
	// Used for easier referencing in other collections and API responses
	Order_id string `json:"order_id"`
	
	// Table_id is the reference to the table where this order was placed (required)
	// This creates a relationship between orders and restaurant tables
	Table_id *string `json:"table_id" validate:"required"`
}
