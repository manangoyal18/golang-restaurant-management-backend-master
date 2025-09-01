package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Note represents a text note in the restaurant management system
// This struct defines the structure of note documents stored in MongoDB
// Notes can be used for special instructions, reminders, or comments
type Note struct {
	// ID is the MongoDB ObjectID - the unique identifier for the note document
	ID primitive.ObjectID `bson:"_id"`
	
	// Text is the main content of the note
	// This contains the detailed note information or message
	Text string `json:"text"`
	
	// Title is the note's heading or subject line
	// This provides a quick summary or identification of the note
	Title string `json:"title"`
	
	// Created_at is the timestamp when the note was created
	Created_at time.Time `json:"created_at"`
	
	// Updated_at is the timestamp when the note was last modified
	Updated_at time.Time `json:"updated_at"`
	
	// Note_id is the string representation of the MongoDB ObjectID
	// Used for easier referencing in other collections and API responses
	Note_id string `json:"note_id"`
}
