package routes

import (
	"project_uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	achievementService *service.AchievementService,
) {
	AuthRoutes(app)
	AchievementRoutes(app, achievementService)
}
