package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// MongoDB connection string
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		// Fallback for local development if not set, though ideally should strictly enforce env
		mongoURI = "mongodb://localhost:27017/ctf"
		log.Println("MONGO_URI not set, using default: " + mongoURI)
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

	// Initialize collections
	Users = DB.Collection("users")
	Challenges = DB.Collection("challenges")
	Submissions = DB.Collection("submissions")
	Teams = DB.Collection("teams")
	Scoreboard = DB.Collection("scoreboard")

	log.Println("Successfully connected to MongoDB!")

	// Create indexes
	createIndexes()
}

func createIndexes() {
	// User indexes
	_, err := Users.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		log.Printf("Error creating user indexes: %v", err)
	}

	// Challenge indexes
	_, err = Challenges.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "category", Value: 1}, {Key: "difficulty", Value: 1}},
		},
	})
	if err != nil {
		log.Printf("Error creating challenge indexes: %v", err)
	}

	// Submission indexes
	_, err = Submissions.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user", Value: 1}, {Key: "challenge", Value: 1}},
		},
	})
	if err != nil {
		log.Printf("Error creating submission indexes: %v", err)
	}

	// Scoreboard index
	_, err = Scoreboard.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "score", Value: -1}, {Key: "lastSolve", Value: 1}},
	})
	if err != nil {
		log.Printf("Error creating scoreboard index: %v", err)
	}
}

// CloseDB closes the MongoDB connection
func CloseDB() {
	if Client != nil {
		Client.Disconnect(context.Background())
	}
}
