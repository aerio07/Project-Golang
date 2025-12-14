package routes

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App) {
	AuthRoutes(app)
	AchievementRoutes(app)
	// UserRoutes(app)
}
