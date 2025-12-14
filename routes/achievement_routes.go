package routes

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(app *fiber.App) {

	// =====================
	// GET achievements
	// =====================
	app.Get(
		"/api/v1/achievements",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:read"),
		service.GetAchievements,
	)

}	
