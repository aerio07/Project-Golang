package main

import (
	"project_uas/app/repository"
	"project_uas/app/service"
	"project_uas/database"
	"project_uas/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.ConnectPostgres()

	achievementRepo := repository.NewAchievementRepository(database.DB)
	achievementService := service.NewAchievementService(achievementRepo)

	app := fiber.New()

	routes.RegisterRoutes(app, achievementService)

	app.Listen(":3000")
}
