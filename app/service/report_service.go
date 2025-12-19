package service

import (
	"project_uas/app/model"
	"project_uas/app/repository"

	"github.com/gofiber/fiber/v2"
)

type ReportService struct {
	StudentRepo  repository.StudentRepository
	LecturerRepo repository.LecturerRepository
	MongoRepo    repository.ReportMongoRepository
}

func NewReportService(stu repository.StudentRepository, lec repository.LecturerRepository, mongo repository.ReportMongoRepository) *ReportService {
	return &ReportService{StudentRepo: stu, LecturerRepo: lec, MongoRepo: mongo}
}

// GetStatistics godoc
// @Summary Get achievement statistics
// @Description Statistik prestasi. Scope otomatis berdasarkan role: Admin=all, Mahasiswa=student, Dosen Wali=advisees
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.ReportStatisticsResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /reports/statistics [get]
func (s *ReportService) GetStatistics(c *fiber.Ctx) error {
	role, userID, ok := getRoleAndUserID(c) // helper kamu (saat ini ada di achievement_service.go)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var filter []string
	scope := "all"

	switch role {
	case "Admin":
		scope = "all"

	case "Mahasiswa":
		studentID, ok2, err := s.StudentRepo.GetByUserID(userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "failed to resolve student"})
		}
		if !ok2 {
			return c.Status(404).JSON(fiber.Map{"message": "student not found"})
		}
		filter = []string{studentID}
		scope = "student"

	case "Dosen Wali":
		lecturerID, ok2, err := s.LecturerRepo.GetByUserID(userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "failed to resolve lecturer"})
		}
		if !ok2 {
			return c.Status(404).JSON(fiber.Map{"message": "lecturer not found"})
		}

		advisees, err := s.LecturerRepo.GetAdvisees(lecturerID, 1000, 0)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"message": "failed to load advisees"})
		}
		for _, a := range advisees {
			filter = append(filter, a.ID)
		}
		scope = "advisees"

	default:
		return c.SendStatus(fiber.StatusForbidden)
	}

	stats, err := s.MongoRepo.AggregateStatistics(filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to generate statistics"})
	}

	if stats == nil {
		stats = &model.AchievementStatistics{}
	}
	stats.Scope = scope
	stats.FilteredStudents = len(filter)

	return c.JSON(fiber.Map{"data": stats})
}

// GetStudentReport godoc
// @Summary Get student report
// @Description Statistik prestasi untuk 1 student (berdasarkan studentId)
// @Tags Reports
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID (uuid)"
// @Success 200 {object} model.ReportStudentResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /reports/students/{id} [get]
func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {
	studentID := c.Params("id")

	stats, err := s.MongoRepo.AggregateStatistics([]string{studentID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to generate student report"})
	}

	if stats == nil {
		stats = &model.AchievementStatistics{}
	}
	stats.Scope = "student"
	stats.FilteredStudents = 1

	return c.JSON(fiber.Map{"data": stats})
}
