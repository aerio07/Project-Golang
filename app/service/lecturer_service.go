package service

import (
	"project_uas/app/repository"

	"github.com/gofiber/fiber/v2"
)

type LecturerService struct {
	Repo repository.LecturerRepository
}

func NewLecturerService(repo repository.LecturerRepository) *LecturerService {
	return &LecturerService{Repo: repo}
}

func (s *LecturerService) GetLecturers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	data, err := s.Repo.List(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch lecturers"})
	}
	return c.JSON(fiber.Map{"data": data})
}

func (s *LecturerService) GetAdvisees(c *fiber.Ctx) error {
	lecturerID := c.Params("id")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	data, err := s.Repo.GetAdvisees(lecturerID, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch advisees"})
	}
	return c.JSON(fiber.Map{"data": data})
}
