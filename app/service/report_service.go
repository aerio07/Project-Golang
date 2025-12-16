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

func (s *ReportService) GetStatistics(c *fiber.Ctx) error {
	role, userID, ok := getRoleAndUserID(c) // ini sudah ada di achievement_service.go
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

func (s *ReportService) GetStudentReport(c *fiber.Ctx) error {
	studentID := c.Params("id")
	stats, err := s.MongoRepo.AggregateStatistics([]string{studentID})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to generate student report"})
	}
	stats.Scope = "student"
	stats.FilteredStudents = 1
	return c.JSON(fiber.Map{"data": stats})
}
