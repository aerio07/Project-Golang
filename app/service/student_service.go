package service

import (
	"strings"

	"project_uas/app/repository"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	Repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) *StudentService {
	return &StudentService{Repo: repo}
}

func (s *StudentService) GetStudents(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	data, err := s.Repo.List(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch students"})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (s *StudentService) GetStudent(c *fiber.Ctx) error {
	id := c.Params("id")
	st, ok, err := s.Repo.GetByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch student"})
	}
	if !ok {
		return c.Status(404).JSON(fiber.Map{"message": "student not found"})
	}
	return c.JSON(fiber.Map{"data": st})
}

func (s *StudentService) GetStudentAchievements(c *fiber.Ctx) error {
	studentID := c.Params("id")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	list, err := s.Repo.GetAchievements(studentID, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch achievements"})
	}
	return c.JSON(fiber.Map{"data": list})
}

type setAdvisorReq struct {
	AdvisorID string `json:"advisor_id"` // lecturers.id (uuid)
}

func (s *StudentService) SetAdvisor(c *fiber.Ctx) error {
	studentID := c.Params("id")

	var body setAdvisorReq
	if err := c.BodyParser(&body); err != nil || strings.TrimSpace(body.AdvisorID) == "" {
		return c.Status(400).JSON(fiber.Map{"message": "advisor_id is required"})
	}

	if err := s.Repo.SetAdvisor(studentID, strings.TrimSpace(body.AdvisorID)); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to set advisor"})
	}
	return c.JSON(fiber.Map{"message": "advisor updated"})
}
