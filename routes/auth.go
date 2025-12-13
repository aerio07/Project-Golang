package routes

import (
	"github.com/gofiber/fiber/v2"
	"project_uas/app/service"
)

func AuthRoutes(app *fiber.App) {
	app.Post("/api/v1/auth/login", service.Login)
}
