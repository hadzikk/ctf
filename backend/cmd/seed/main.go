package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"ctf-backend/database"
	"ctf-backend/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Try loading .env from probable locations, ignore errors as InitDB also checks
	_ = godotenv.Load()             // Check current directory
	_ = godotenv.Load("../../.env") // Check root if running from cmd/seed

	database.InitDB()
	defer database.CloseDB()

	seedUsers()
	seedChallenges()

	fmt.Println("Seeding completed successfully!")
}

func seedUsers() {
	ctx := context.Background()

	// 1. Admin User
	adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Username:   "admin",
		Email:      "admin@ctf.com",
		Password:   string(adminPass),
		Role:       "admin",
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
	}

	// Upsert Admin
	_, err := database.Users.UpdateOne(ctx,
		bson.M{"username": admin.Username},
		bson.M{"$set": admin},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to seed admin: %v", err)
	} else {
		log.Println("Seeded admin user")
	}

	// 2. Standard User
	userPass, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	user := models.User{
		Username:   "player1",
		Email:      "player1@ctf.com",
		Password:   string(userPass),
		Role:       "user",
		CreatedAt:  time.Now(),
		LastActive: time.Now(),
	}

	_, err = database.Users.UpdateOne(ctx,
		bson.M{"username": user.Username},
		bson.M{"$set": user},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Printf("Failed to seed user: %v", err)
	} else {
		log.Println("Seeded standard user")
	}
}

// Helper to get pointer options for Upsert
// Since go.mongodb.org/mongo-driver/mongo/options is not imported directly in this helper yet
// We can use the one from database package if exported, but it's not.
// Let's just import options above.
// Wait, I forgot to import options. Adding it now.

func seedChallenges() {
	ctx := context.Background()

	challenges := []models.Challenge{
		{
			Title:       "GPS Spoofing",
			Category:    "GIS",
			Difficulty:  "Easy",
			Points:      100,
			Description: "You need to find the hidden location by analyzing the GPS coordinates in the request.",
			Flag:        "CTF{gps_sp00f1ng_1s_fun}",
			AuthorID:    primitive.NewObjectID(), // Mock ID
			IsActive:    true,
			MapConfig: &models.MapConfig{
				Center: models.MapPoint{Lat: -6.2088, Lng: 106.8456}, // Jakarta
				Zoom:   13,
				Markers: []models.MapMarker{
					{
						Position:    models.MapPoint{Lat: -6.2088, Lng: 106.8456},
						Title:       "Monas",
						Description: "The National Monument",
					},
				},
			},
		},
		{
			Title:       "GeoJSON Injection",
			Category:    "Web",
			Difficulty:  "Medium",
			Points:      200,
			Description: "Exploit a vulnerability in the GeoJSON parsing to reveal the flag.",
			Flag:        "CTF{ge0j50n_1nj3ct10n_ftw}",
			AuthorID:    primitive.NewObjectID(),
			IsActive:    true,
			MapConfig: &models.MapConfig{
				Center: models.MapPoint{Lat: -7.2575, Lng: 112.7521}, // Surabaya
				Zoom:   12,
			},
		},
		{
			Title:       "Geofence Escape",
			Category:    "GIS",
			Difficulty:  "Hard",
			Points:      300,
			Description: "The application has a geofence that restricts certain actions. Find a way to bypass this restriction.",
			Flag:        "CTF{g30f3nc3_byp4ss3d}",
			AuthorID:    primitive.NewObjectID(),
			IsActive:    true,
			MapConfig: &models.MapConfig{
				Center: models.MapPoint{Lat: -6.9147, Lng: 107.6098}, // Bandung
				Zoom:   12,
			},
		},
	}

	for _, ch := range challenges {
		// Prepare document (ensure IDs are handled if needed, or let Mongo generate)
		// We use Title as unique key for seeding idempotency
		_, err := database.Challenges.UpdateOne(ctx,
			bson.M{"title": ch.Title},
			bson.M{"$set": ch},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			log.Printf("Failed to seed challenge %s: %v", ch.Title, err)
		} else {
			log.Printf("Seeded challenge: %s", ch.Title)
		}
	}
}
