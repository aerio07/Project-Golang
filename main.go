package main

import (
	"project_uas/database"
	"project_uas/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.ConnectPostgres()

	app := fiber.New()

	// REGISTER ALL ROUTES
	routes.RegisterRoutes(app)

	app.Listen(":3000")
}
