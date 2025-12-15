package routes

import (
	"project_uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	authService *service.AuthService,
	achievementService *service.AchievementService,
) {
	AuthRoutes(app, authService)
	AchievementRoutes(app, achievementService)
}

