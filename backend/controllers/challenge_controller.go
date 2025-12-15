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

type ChallengeController struct {
	collection *mongo.Collection
}

func NewChallengeController(db *mongo.Database) *ChallengeController {
	return &ChallengeController{
		collection: db.Collection("challenges"),
	}
}

func (cc *ChallengeController) GetAllChallenges(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var challenges []models.Challenge
	cursor, err := cc.collection.Find(ctx, bson.M{"isActive": true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch challenges",
		})
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &challenges); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode challenges",
		})
	}

	return c.JSON(challenges)
}

func (cc *ChallengeController) GetChallengeByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid challenge ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var challenge models.Challenge
	err = cc.collection.FindOne(ctx, bson.M{"_id": objID, "isActive": true}).Decode(&challenge)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Challenge not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch challenge",
		})
	}

	return c.JSON(challenge)
}

func (cc *ChallengeController) CreateChallenge(c *fiber.Ctx) error {
	var challenge models.Challenge

	if err := c.BodyParser(&challenge); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Set default values
	challenge.IsActive = true
	challenge.CreatedAt = time.Now()
	challenge.UpdatedAt = time.Now()

	// Insert into database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := cc.collection.InsertOne(ctx, challenge)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create challenge",
		})
	}

	challenge.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(challenge)
}

func (cc *ChallengeController) UpdateChallenge(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid challenge ID",
		})
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Update the updatedAt field
	updateData["updatedAt"] = time.Now()

	// Update in database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := cc.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": updateData},
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update challenge",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Challenge not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Challenge updated successfully",
	})
}

func (cc *ChallengeController) DeleteChallenge(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid challenge ID",
		})
	}

	// Soft delete by setting isActive to false
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := cc.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"isActive": false, "updatedAt": time.Now()}},
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete challenge",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Challenge not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Challenge deleted successfully",
	})
}
