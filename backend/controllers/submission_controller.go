package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"ctf-backend/models"
)

type SubmissionController struct {
	submissionCollection *mongo.Collection
	challengeCollection  *mongo.Collection
}

func NewSubmissionController(db *mongo.Database) *SubmissionController {
	return &SubmissionController{
		submissionCollection: db.Collection("submissions"),
		challengeCollection:  db.Collection("challenges"),
	}
}

func (sc *SubmissionController) SubmitFlag(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjID, _ := primitive.ObjectIDFromHex(userID)

	var submission struct {
		ChallengeID string `json:"challengeId"`
		Flag        string `json:"flag"`
	}

	if err := c.BodyParser(&submission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	challengeID, err := primitive.ObjectIDFromHex(submission.ChallengeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid challenge ID",
		})
	}

	// Check if challenge exists and get the correct flag
	var challenge models.Challenge
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = sc.challengeCollection.FindOne(ctx, bson.M{
		"_id":      challengeID,
		"isActive": true,
	}).Decode(&challenge)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Challenge not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify flag",
		})
	}

	// Check if already solved
	var existingSubmission models.Submission
	err = sc.submissionCollection.FindOne(ctx, bson.M{
		"userID":      userObjID,
		"challengeID": challengeID,
		"isCorrect":   true,
	}).Decode(&existingSubmission)

	if err == nil {
		return c.JSON(fiber.Map{
			"correct": true,
			"message": "You've already solved this challenge!",
			"points":  challenge.Points,
		})
	}

	// Check if flag is correct
	isCorrect := submission.Flag == challenge.Flag

	// Save submission
	submissionDoc := models.Submission{
		ID:            primitive.NewObjectID(),
		UserID:        userObjID,
		ChallengeID:   challengeID,
		Flag:          submission.Flag,
		IsCorrect:     isCorrect,
		PointsAwarded: 0, // Will be set based on isCorrect
	}

	if isCorrect {
		submissionDoc.PointsAwarded = challenge.Points
	}

	submissionDoc.BeforeCreate() // Set CreatedAt

	_, err = sc.submissionCollection.InsertOne(ctx, submissionDoc)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save submission",
		})
	}

	if isCorrect {
		return c.JSON(fiber.Map{
			"correct": true,
			"message": "Correct flag! Well done!",
			"points":  challenge.Points,
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"correct": false,
		"message": "Incorrect flag. Try again!",
	})
}

func (sc *SubmissionController) GetUserSubmissions(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjID, _ := primitive.ObjectIDFromHex(userID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get user's submissions with challenge details
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"userID": userObjID}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "challenges",
			"localField":   "challengeID",
			"foreignField": "_id",
			"as":           "challenge",
		}}},
		bson.D{{Key: "$unwind", Value: "$challenge"}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id":            1,
			"challengeID":    1,
			"challengeTitle": "$challenge.title",
			"category":       "$challenge.category",
			"points":         1,
			"isCorrect":      1,
			"submittedAt":    1,
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"submittedAt": -1}}}} // Newest first

	cursor, err := sc.submissionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch submissions",
		})
	}
	defer cursor.Close(ctx)

	var submissions []bson.M
	if err = cursor.All(ctx, &submissions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode submissions",
		})
	}

	return c.JSON(submissions)
}

func (sc *SubmissionController) GetAllSubmissions(c *fiber.Ctx) error {
	// Only show correct submissions for all users (admin view)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"isCorrect": true}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "challenges",
			"localField":   "challengeID",
			"foreignField": "_id",
			"as":           "challenge",
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "users",
			"localField":   "userID",
			"foreignField": "_id",
			"as":           "user",
		}}},
		bson.D{{Key: "$unwind", Value: "$challenge"}},
		bson.D{{Key: "$unwind", Value: "$user"}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id":         1,
			"username":    "$user.username",
			"challenge":   "$challenge.title",
			"category":    "$challenge.category",
			"points":      "$challenge.points",
			"submittedAt": 1,
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"submittedAt": -1}}}} // Newest first

	cursor, err := sc.submissionCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch submissions",
		})
	}
	defer cursor.Close(ctx)

	var submissions []bson.M
	if err = cursor.All(ctx, &submissions); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode submissions",
		})
	}

	return c.JSON(submissions)
}
