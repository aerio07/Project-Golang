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

// GetLecturers godoc
// @Summary List lecturers
// @Description Ambil daftar dosen (pagination)
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} model.LecturerListResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /lecturers [get]
func (s *LecturerService) GetLecturers(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	data, err := s.Repo.List(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch lecturers"})
	}
	return c.JSON(fiber.Map{"data": data})
}

// GetAdvisees godoc
// @Summary List advisees
// @Description Ambil daftar mahasiswa bimbingan dari dosen (by lecturer id)
// @Tags Lecturers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Lecturer ID (uuid)"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} model.AdviseeListResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /lecturers/{id}/advisees [get]
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
