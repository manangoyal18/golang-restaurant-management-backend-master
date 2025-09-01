package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Table represents a restaurant table in the management system
// This struct defines the structure of table documents stored in MongoDB
// Each table has capacity information and unique identification
type Table struct {
	// ID is the MongoDB ObjectID - the unique identifier for the table document
	ID primitive.ObjectID `bson:"_id"`
	
	// Number_of_guests is the current number of guests seated at this table (required)
	// This helps with capacity management and service planning
	Number_of_guests *int `json:"number_of_guests" validate:"required"`
	
	// Table_number is the unique table identifier used by restaurant staff (required)
	// This is the physical table number displayed in the restaurant
	Table_number *int `json:"table_number" validate:"required"`
	
	// Created_at is the timestamp when the table record was created
	Created_at time.Time `json:"created_at"`
	
	// Updated_at is the timestamp when the table information was last modified
	Updated_at time.Time `json:"updated_at"`
	
	// Table_id is the string representation of the MongoDB ObjectID
	// Used for easier referencing in other collections and API responses
	Table_id string `json:"table_id"`
}
