package routes

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(app *fiber.App) {

	// GET achievements
	app.Get(
		"/api/v1/achievements",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:read"),
		service.GetAchievements,
	)

	// POST achievements (submit) - Mahasiswa
	app.Post(
		"/api/v1/achievements",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:create"),
		service.CreateAchievement,
	)

	// VERIFY achievement (Dosen Wali)
	app.Post(
		"/api/v1/achievements/:id/verify",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:verify"),
		service.VerifyAchievement,
	)

	// REJECT achievement (Dosen Wali)
	app.Post(
		"/api/v1/achievements/:id/reject",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:verify"),
		service.RejectAchievement,
	)
}
