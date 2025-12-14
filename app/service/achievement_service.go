	package service

	import (
		"project_uas/app/repository"

		"github.com/gofiber/fiber/v2"
		"github.com/golang-jwt/jwt/v5"
	)

	func GetAchievements(c *fiber.Ctx) error {
		// ambil token dari context
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)

		userID := claims["sub"].(string)
		role := claims["role"].(string)

		var data interface{}
		var err error

		switch role {
		case "Admin":
			data, err = repository.GetAllAchievements()
		case "Mahasiswa":
			data, err = repository.GetAchievementsByStudent(userID)
		case "Dosen Wali":
			data, err = repository.GetAchievementsBySupervisor(userID)
		default:
			return c.Status(403).JSON(fiber.Map{
				"message": "role not allowed",
			})
		}

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"message": "failed to fetch achievements",
			})
		}

		return c.JSON(fiber.Map{
			"data": data,
		})
	}
