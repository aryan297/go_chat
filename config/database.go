package config

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

// InitMongoDB initializes the MongoDB connection
func InitMongoDB() {
	var err error

	// Connect to MongoDB
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection to ensure it's valid
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB!")
}

// GetDB returns the MongoDB database instance
func GetDB() *mongo.Database {
	if client == nil {
		log.Fatal("MongoDB client is not initialized. Please call InitMongoDB first.")
	}
	return client.Database("chat-app")
}

// CloseMongoDB gracefully closes the MongoDB connection
func CloseMongoDB() {
	if client != nil {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
		log.Println("Disconnected from MongoDB")
	}
}
