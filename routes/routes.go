package routes

import (
	"project_uas/app/service"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	authService *service.AuthService,
	achievementService *service.AchievementService,
	userService *service.UserService,
	studentService *service.StudentService,
	lecturerService *service.LecturerService,
	reportService *service.ReportService,
) {
	AuthRoutes(app, authService)
	AchievementRoutes(app, achievementService)

	UserRoutes(app, userService)
	StudentRoutes(app, studentService)
	LecturerRoutes(app, lecturerService)
	ReportRoutes(app, reportService)
}
