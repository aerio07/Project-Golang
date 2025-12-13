package main

import (
	"github.com/gofiber/fiber/v2"
	"project_uas/database"
	"project_uas/routes"
)

func main() {
	database.ConnectPostgres()

	app := fiber.New()

	routes.AuthRoutes(app)

	app.Listen(":3000")
}
