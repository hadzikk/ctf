package routes

import (
	"ctf-backend/controllers"
	"ctf-backend/database"
	"ctf-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupSubmissionRoutes(api fiber.Router) {
	// Use global database instance
	submissionController := controllers.NewSubmissionController(database.DB)

	submissionRoutes := api.Group("/submissions")
	submissionRoutes.Use(middleware.RequireAuth())
	{
		submissionRoutes.Post("/", submissionController.SubmitFlag)
		submissionRoutes.Get("/", submissionController.GetUserSubmissions)

		// Admin only
		adminRoutes := submissionRoutes.Group("/admin")
		adminRoutes.Use(middleware.RequireAdmin())
		{
			adminRoutes.Get("/all", submissionController.GetAllSubmissions)
		}
	}
}
