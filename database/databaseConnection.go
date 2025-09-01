// Package database handles all MongoDB database connections and collection management
// This package provides utilities to connect to MongoDB and access specific collections
package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBinstance creates and returns a MongoDB client connection
// This function establishes a connection to the local MongoDB instance
// Returns: *mongo.Client - A connected MongoDB client instance
func DBinstance() *mongo.Client {
	// MongoDB connection string - points to local MongoDB instance on default port
	// In production, this should be moved to environment variables for security
	MongoDb := "mongodb://localhost:27017"
	fmt.Print(MongoDb)

	// Create a new MongoDB client with the connection string
	// This prepares the client but doesn't establish the connection yet
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}
	
	// Create a context with a 10-second timeout for the connection attempt
	// This prevents the application from hanging indefinitely if MongoDB is unavailable
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Ensure the context is cancelled to free resources
	defer cancel()

	// Actually establish the connection to MongoDB
	// This is when the real network connection is made
	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to mongodb")
	return client
}

// Client is a global MongoDB client instance that's initialized at package load time
// This provides a singleton pattern for database access across the application
var Client *mongo.Client = DBinstance()

// OpenCollection returns a reference to a specific MongoDB collection
// Parameters:
//   - client: The MongoDB client instance
//   - collectionName: The name of the collection to access
// Returns: *mongo.Collection - A reference to the specified collection in the "restaurant" database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	// Access the "restaurant" database and the specified collection
	// All collections in this application are stored under the "restaurant" database
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)

	return collection
}
