package routes

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func ReportRoutes(app *fiber.App, svc *service.ReportService) {
	base := "/api/v1/reports"

	app.Get(base+"/statistics",
		middleware.JWTMiddleware,
		middleware.RequirePermission("report:read"),
		svc.GetStatistics,
	)

	app.Get(base+"/student/:id",
		middleware.JWTMiddleware,
		middleware.RequirePermission("report:read"),
		svc.GetStudentReport,
	)
}
