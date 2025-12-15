package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Initialize API v1 routes
	api := app.Group("/api/v1")

	// Setup routes for each domain
	SetupUserRoutes(api)
	SetupChallengeRoutes(api)
	SetupSubmissionRoutes(api)
	SetupTeamRoutes(api)
}
