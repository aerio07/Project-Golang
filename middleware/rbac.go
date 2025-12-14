package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		perms, ok := claims["permissions"].([]interface{})
		if !ok {
			return c.Status(403).JSON(fiber.Map{
				"message": "permissions not found",
			})
		}

		for _, p := range perms {
			if p.(string) == permission {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"message": "forbidden",
		})
	}
}
