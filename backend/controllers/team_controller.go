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

type TeamController struct {
	collection     *mongo.Collection
	userCollection *mongo.Collection
}

func NewTeamController(db *mongo.Database) *TeamController {
	return &TeamController{
		collection:     db.Collection("teams"),
		userCollection: db.Collection("users"),
	}
}

func (tc *TeamController) CreateTeam(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	userObjID, _ := primitive.ObjectIDFromHex(userID)

	var team models.Team
	if err := c.BodyParser(&team); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Check if user is already in a team
	// (Simplification: assuming user model has team_id or we check team members)
	// For now, let's just create the team

	team.CaptainID = userObjID
	team.Members = []primitive.ObjectID{userObjID}
	team.BeforeCreate()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := tc.collection.InsertOne(ctx, team)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create team. Name might be taken.",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Team created successfully",
		"teamId":  result.InsertedID,
	})
}

func (tc *TeamController) GetAllTeams(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var teams []models.Team
	cursor, err := tc.collection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch teams",
		})
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &teams); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to decode teams",
		})
	}

	return c.JSON(teams)
}

func (tc *TeamController) GetTeamByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid team ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var team models.Team
	err = tc.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&team)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Team not found",
		})
	}

	return c.JSON(team)
}

func (tc *TeamController) UpdateTeam(c *fiber.Ctx) error {
	// TODO: Implement update logic (e.g. rename)
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (tc *TeamController) DeleteTeam(c *fiber.Ctx) error {
	// TODO: Implement delete logic
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (tc *TeamController) JoinTeam(c *fiber.Ctx) error {
	// TODO: Implement join logic
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (tc *TeamController) LeaveTeam(c *fiber.Ctx) error {
	// TODO: Implement leave logic
	return c.SendStatus(fiber.StatusNotImplemented)
}
