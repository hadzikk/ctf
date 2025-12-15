package routes

import (
	"ctf-backend/controllers"
	"ctf-backend/database"
	"ctf-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupTeamRoutes(api fiber.Router) {
	// Use global database instance
	teamController := controllers.NewTeamController(database.DB)

	teamRoutes := api.Group("/teams")
	{
		// Public routes
		teamRoutes.Get("/", teamController.GetAllTeams)
		teamRoutes.Get("/:id", teamController.GetTeamByID)

		// Protected routes
		teamRoutes.Use(middleware.RequireAuth())
		teamRoutes.Post("/", teamController.CreateTeam)
		teamRoutes.Put("/:id", teamController.UpdateTeam)
		teamRoutes.Delete("/:id", teamController.DeleteTeam)
		teamRoutes.Post("/:id/join", teamController.JoinTeam)
		teamRoutes.Post("/:id/leave", teamController.LeaveTeam)
	}
}
