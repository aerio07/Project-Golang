package routes

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func AchievementRoutes(app *fiber.App, svc *service.AchievementService) {

	base := "/api/v1/achievements"

	app.Get(base,
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:read"),
		svc.GetAchievements,
	)

	app.Get(base+"/:id",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:read"),
		svc.GetAchievementDetail,
	)

	app.Post(base,
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:create"),
		svc.CreateAchievement,
	)

	app.Put(base+"/:id",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:update"),
		svc.UpdateAchievement,
	)

	app.Delete(base+"/:id",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:delete"),
		svc.DeleteAchievement,
	)

	app.Post(base+"/:id/submit",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:update"),
		svc.SubmitAchievement,
	)

	app.Post(base+"/:id/verify",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:verify"),
		svc.VerifyAchievement,
	)

	app.Post(base+"/:id/reject",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:verify"),
		svc.RejectAchievement,
	)

	app.Get(base+"/:id/history",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:read"),
		svc.GetAchievementHistory,
	)

	app.Post(base+"/:id/attachments",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:update"),
		svc.UploadAchievementAttachment,
	)
}
