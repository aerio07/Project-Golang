package routes

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func StudentRoutes(app *fiber.App, svc *service.StudentService) {
	base := "/api/v1/students"

	app.Get(base,
		middleware.JWTMiddleware,
		middleware.RequirePermission("student:read"),
		svc.GetStudents,
	)

	app.Get(base+"/:id",
		middleware.JWTMiddleware,
		middleware.RequirePermission("student:read"),
		svc.GetStudent,
	)

	app.Get(base+"/:id/achievements",
		middleware.JWTMiddleware,
		middleware.RequirePermission("achievement:read"),
		svc.GetStudentAchievements,
	)

	app.Put(base+"/:id/advisor",
		middleware.JWTMiddleware,
		middleware.RequirePermission("student:update"),
		svc.SetAdvisor,
	)
}
