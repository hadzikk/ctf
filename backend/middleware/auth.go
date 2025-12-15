package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireAuth is a middleware to check if the request is authenticated
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Missing token",
			})
		}

		// TODO: Validate JWT token
		// For now, we'll just check if the token exists

		// Add user to context if needed
		// c.Locals("user", user)

		return c.Next()
	}
}
