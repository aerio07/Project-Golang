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
	database.ConnectMongo()

	achievementRepo := repository.NewAchievementRepository(database.DB)
	achievementMongoRepo := repository.NewAchievementMongoRepository(database.MongoDB)
	achievementService := service.NewAchievementService(achievementRepo, achievementMongoRepo)

	authRepo := repository.NewAuthRepository(database.DB)
	authService := service.NewAuthService(authRepo)

	app := fiber.New()

	// serve uploads
	app.Static("/uploads", "./uploads")

	routes.RegisterRoutes(app, authService, achievementService)

	app.Listen(":3000")
}
