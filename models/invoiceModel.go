package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Invoice represents a billing invoice in the restaurant management system
// This struct defines the structure of invoice documents stored in MongoDB
// Each invoice is associated with an order and tracks payment information
type Invoice struct {
	// ID is the MongoDB ObjectID - the unique identifier for the invoice document
	ID primitive.ObjectID `bson:"_id"`
	
	// Invoice_id is the string representation of the MongoDB ObjectID
	// Used for easier referencing in other collections and API responses
	Invoice_id string `json:"invoice_id"`
	
	// Order_id is the reference to the order this invoice is for
	// This creates a relationship between invoices and orders
	Order_id string `json:"order_id"`
	
	// Payment_method is how the customer will pay (CARD, CASH, or empty for not specified)
	// The validation ensures only valid payment methods are accepted
	Payment_method *string `json:"payment_method" validate:"eq=CARD|eq=CASH|eq="`
	
	// Payment_status tracks whether the invoice has been paid (required: PENDING or PAID)
	// This is used for financial tracking and order completion
	Payment_status *string `json:"payment_status" validate:"required,eq=PENDING|eq=PAID"`
	
	// Payment_due_date is when the payment is due
	// Used for tracking overdue payments and follow-up
	Payment_due_date time.Time `json:"Payment_due_date"`
	
	// Created_at is the timestamp when the invoice was generated
	Created_at time.Time `json:"created_at"`
	
	// Updated_at is the timestamp when the invoice was last modified
	Updated_at time.Time `json:"updated_at"`
}
