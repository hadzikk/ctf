package routes

import (
	"ctf-backend/controllers"
	"ctf-backend/database"
	"ctf-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupChallengeRoutes(api fiber.Router) {
	// Use global database instance
	challengeController := controllers.NewChallengeController(database.DB)

	challengeRoutes := api.Group("/challenges")
	{
		// Public routes
		challengeRoutes.Get("/", challengeController.GetAllChallenges)
		challengeRoutes.Get("/:id", challengeController.GetChallengeByID)

		// Protected routes (require admin)
		challengeRoutes.Use(middleware.RequireAuth(), middleware.RequireAdmin())
		challengeRoutes.Post("/", challengeController.CreateChallenge)
		challengeRoutes.Put("/:id", challengeController.UpdateChallenge)
		challengeRoutes.Delete("/:id", challengeController.DeleteChallenge)
	}
}
