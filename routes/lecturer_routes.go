package routes

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func LecturerRoutes(app *fiber.App, svc *service.LecturerService) {
	base := "/api/v1/lecturers"

	app.Get(base,
		middleware.JWTMiddleware,
		middleware.RequirePermission("lecturer:read"),
		svc.GetLecturers,
	)

	app.Get(base+"/:id/advisees",
		middleware.JWTMiddleware,
		middleware.RequirePermission("student:read"),
		svc.GetAdvisees,
	)
}
