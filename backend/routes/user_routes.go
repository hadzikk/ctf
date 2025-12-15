package routes

import (
	"ctf-backend/controllers"
	"ctf-backend/database"
	"ctf-backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router) {
	// Use global database instance
	userController := controllers.NewUserController(database.DB)

	userRoutes := api.Group("/users")
	{
		userRoutes.Post("/register", userController.Register)
		userRoutes.Post("/login", userController.Login)

		// Protected routes (require authentication)
		userRoutes.Use(middleware.RequireAuth())
		userRoutes.Get("/me", userController.GetCurrentUser)
		userRoutes.Put("/me", userController.UpdateProfile)
	}
}
