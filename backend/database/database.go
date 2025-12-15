package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client      *mongo.Client
	DB          *mongo.Database
	Users       *mongo.Collection
	Challenges  *mongo.Collection
	Submissions *mongo.Collection
	Teams       *mongo.Collection
	Scoreboard  *mongo.Collection
)

func InitDB() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// MongoDB connection string
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI must be set in .env file")
	}

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	Client = client
	DB = client.Database("ctf")

	// Initialize database and collections
	DB = client.Database("ctf")
	Users = DB.Collection("users")
	Challenges = DB.Collection("challenges")
	Submissions = DB.Collection("submissions")
	Teams = DB.Collection("teams")
	Scoreboard = DB.Collection("scoreboard")

	log.Println("Successfully connected to MongoDB Atlas!")

	// Create indexes
	createIndexes()
}

func createIndexes() {
	// User indexes
	_, err := Users.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"username": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		log.Fatalf("Error creating user indexes: %v", err)
	}

	// Challenge indexes
	_, err = Challenges.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"category": 1, "difficulty": 1},
		},
	})
	if err != nil {
		log.Fatalf("Error creating challenge indexes: %v", err)
	}

	// Submission indexes
	_, err = Submissions.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: map[string]interface{}{"user": 1, "challenge": 1},
		},
	})
	if err != nil {
		log.Fatalf("Error creating submission indexes: %v", err)
	}

	// Scoreboard index
	_, err = Scoreboard.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: map[string]interface{}{"score": -1, "lastSolve": 1},
	})
	if err != nil {
		log.Fatalf("Error creating scoreboard index: %v", err)
	}
}

// CloseDB closes the MongoDB connection
func CloseDB() {
	if Client != nil {
		Client.Disconnect(context.Background())
	}
}
