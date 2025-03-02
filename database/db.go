package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection

func ConnectDB() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load the .env file")
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping the database
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	fmt.Println("✅ Connected to MongoDB")

	// Set global collection variable
	Collection = client.Database("todo-go-lang").Collection("todo")
}
