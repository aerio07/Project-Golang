package routes

import (
	"project_uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authService *service.AuthService) {
	app.Post("/api/v1/auth/login", authService.Login)
}
