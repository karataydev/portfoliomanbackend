package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JwtAuthMiddleware(tokenService *TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization token",
			})
		}

		// Remove 'Bearer ' prefix if present
		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := tokenService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Set user information in context for use in handlers
		c.Locals("userId", claims.UserId)
		c.Locals("userEmail", claims.Email)

		return c.Next()
	}
}
