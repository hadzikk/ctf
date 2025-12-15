//go:build test
// +build test

package database_test

import (
	"log"
	"os"
	"testing"

	"ctf-backend/database"
)

func TestDatabaseConnection(t *testing.T) {
	// Change to project root directory
	err := os.Chdir("..")
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Test connection by pinging the database
	err = database.Client.Ping(nil, nil)
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Successfully connected to MongoDB Atlas!")
}
