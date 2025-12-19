package main

import (
	"project_uas/app/repository"
	"project_uas/app/service"
	"project_uas/database"
	"project_uas/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	// Swagger
	_ "project_uas/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	
)

// @title Sistem Pelaporan Prestasi Mahasiswa API
// @version 1.0
// @description Dokumentasi API Sistem Pelaporan Prestasi Mahasiswa
// @host localhost:3000
// @BasePath /api/v1
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	_ = godotenv.Load()
	
	app := fiber.New()
	
	// ✅ CORS biar Swagger UI bisa call API
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	
	// ✅ Swagger UI
	app.Get("/swagger/*", fiberSwagger.WrapHandler)




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
