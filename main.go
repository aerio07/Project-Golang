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
	app := fiber.New()
	
	database.ConnectPostgres()
	database.ConnectMongo()

	achievementRepo := repository.NewAchievementRepository(database.DB)
	achievementMongoRepo := repository.NewAchievementMongoRepository(database.MongoDB)
	achievementService := service.NewAchievementService(achievementRepo, achievementMongoRepo)

	authRepo := repository.NewAuthRepository(database.DB)
	authService := service.NewAuthService(authRepo)

	userRepo := repository.NewUserRepository(database.DB)
userService := service.NewUserService(userRepo)

studentRepo := repository.NewStudentRepository(database.DB)
studentService := service.NewStudentService(studentRepo)

lecturerRepo := repository.NewLecturerRepository(database.DB)
lecturerService := service.NewLecturerService(lecturerRepo)

reportMongoRepo := repository.NewReportMongoRepository(database.MongoDB)
reportService := service.NewReportService(studentRepo, lecturerRepo, reportMongoRepo)

routes.RegisterRoutes(app, authService, achievementService, userService, studentService, lecturerService, reportService)



	// serve uploads
	app.Static("/uploads", "./uploads")


	app.Listen(":3000")
}
