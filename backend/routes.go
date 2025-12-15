package main

import (
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// User routes
	userRoutes := api.Group("/users")
	{
		userRoutes.Post("/register", registerUser)
		userRoutes.Post("/login", loginUser)
		userRoutes.Get("/me", authenticateUser, getCurrentUser)
	}

	// Challenge routes
	challengeRoutes := api.Group("/challenges")
	{
		challengeRoutes.Get("/", listChallenges)
		challengeRoutes.Get("/:id", getChallenge)
		challengeRoutes.Post("/", authenticateUser, createChallenge)      // Admin only
		challengeRoutes.Put("/:id", authenticateUser, updateChallenge)    // Admin only
		challengeRoutes.Delete("/:id", authenticateUser, deleteChallenge) // Admin only
	}

	// Submission routes
	submissionRoutes := api.Group("/submissions")
	submissionRoutes.Use(authenticateUser)
	{
		submissionRoutes.Post("/", submitFlag)
		submissionRoutes.Get("/me", getUserSubmissions)
		submissionRoutes.Get("/:id", getSubmission)
	}

	// Team routes
	teamRoutes := api.Group("/teams")
	teamRoutes.Use(authenticateUser)
	{
		teamRoutes.Post("/", createTeam)
		teamRoutes.Get("/me", getMyTeam)
		teamRoutes.Post("/join", joinTeam)
		teamRoutes.Post("/leave", leaveTeam)
	}
}

// Authentication middleware
func authenticateUser(c *fiber.Ctx) error {
	// TODO: Implement JWT authentication
	// For now, just a placeholder
	return c.Next()
}
