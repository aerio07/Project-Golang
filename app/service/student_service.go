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

// GetStudents godoc
// @Summary List students
// @Description Ambil daftar mahasiswa (pagination)
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} model.StudentListResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /students [get]
func (s *StudentService) GetStudents(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	data, err := s.Repo.List(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch students"})
	}
	return c.JSON(fiber.Map{"data": data})
}

// GetStudent godoc
// @Summary Get student detail
// @Description Ambil detail mahasiswa by id
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID (uuid)"
// @Success 200 {object} model.StudentDetailResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /students/{id} [get]
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

// GetStudentAchievements godoc
// @Summary List student achievements
// @Description Ambil list prestasi milik student by studentId (uuid)
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param id path string true "Student ID (uuid)"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} model.StudentAchievementListResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /students/{id}/achievements [get]
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

// SetAdvisor godoc
// @Summary Set student advisor
// @Description Set dosen wali untuk mahasiswa (advisor_id = lecturers.id)
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Student ID (uuid)"
// @Param body body model.StudentSetAdvisorRequest true "Request body"
// @Success 200 {object} model.MessageResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /students/{id}/advisor [put]
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
