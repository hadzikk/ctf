package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/ctf/backend/controllers"
)

func SetupUserRoutes(api fiber.Router) {
	userController := controllers.NewUserController()

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
