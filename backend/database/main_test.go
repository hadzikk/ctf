package database

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	// Load .env file
	err := os.Chdir("..")
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Initialize database
	InitDB()
	defer CloseDB()

	// Test connection by pinging the database
	err = Client.Ping(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to MongoDB Atlas!")
}
