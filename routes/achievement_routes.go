package routes

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(app *fiber.App, achievementService *service.AchievementService) {

	app.Get(
		"/api/v1/achievements",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:read"),
		achievementService.GetAchievements,
	)

	app.Post(
		"/api/v1/achievements",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:create"),
		achievementService.CreateAchievement,
	)

	app.Post(
		"/api/v1/achievements/:id/submit",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:create"),
		achievementService.SubmitAchievement,
	)

	app.Delete(
		"/api/v1/achievements/:id",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:delete"),
		achievementService.DeleteAchievement,
	)

	app.Post(
		"/api/v1/achievements/:id/verify",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:verify"),
		achievementService.VerifyAchievement,
	)

	app.Post(
		"/api/v1/achievements/:id/reject",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:verify"),
		achievementService.RejectAchievement,
	)
}
