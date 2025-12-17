package routes

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authService *service.AuthService) {
	base := "/api/v1/auth"

	app.Post(base+"/login", authService.Login)
	app.Post(base+"/refresh", authService.Refresh)

	app.Post(base+"/logout",
		middleware.JWTMiddleware,
		authService.Logout,
	)

	app.Get(base+"/profile",
		middleware.JWTMiddleware,
		authService.Profile,
	)
}
