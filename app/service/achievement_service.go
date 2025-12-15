package service

import (
	"project_uas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AchievementService struct {
	Repo repository.AchievementRepository
}

func NewAchievementService(repo repository.AchievementRepository) *AchievementService {
	return &AchievementService{Repo: repo}
}

// ===================== GET =====================

func (s *AchievementService) GetAchievements(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	role := claims["role"].(string)
	userID := claims["sub"].(string)

	var (
		data interface{}
		err  error
	)

	switch role {
	case "Admin":
		data, err = s.Repo.GetAll()
	case "Mahasiswa":
		data, err = s.Repo.GetByStudent(userID)
	case "Dosen Wali":
		data, err = s.Repo.GetBySupervisor(userID)
	default:
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch achievements"})
	}

	return c.JSON(fiber.Map{"data": data})
}

// ===================== CREATE =====================

func (s *AchievementService) CreateAchievement(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	if token.Claims.(jwt.MapClaims)["role"] != "Mahasiswa" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	var body struct {
		StudentID string `json:"student_id"`
	}
	if err := c.BodyParser(&body); err != nil || body.StudentID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "student_id is required"})
	}

	if err := s.Repo.CreateDraft(body.StudentID); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to create achievement"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "achievement created (draft)"})
}

// ===================== SUBMIT =====================

func (s *AchievementService) SubmitAchievement(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	if token.Claims.(jwt.MapClaims)["role"] != "Mahasiswa" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err := s.Repo.Submit(c.Params("id"), token.Claims.(jwt.MapClaims)["sub"].(string)); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to submit achievement"})
	}

	return c.JSON(fiber.Map{"message": "achievement submitted"})
}

// ===================== DELETE =====================

func (s *AchievementService) DeleteAchievement(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	if claims["role"] != "Mahasiswa" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	ok, err := s.Repo.CanDelete(c.Params("id"), claims["sub"].(string))
	if err != nil || !ok {
		return c.Status(403).JSON(fiber.Map{"message": "achievement cannot be deleted"})
	}

	if err := s.Repo.SoftDelete(c.Params("id")); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to delete achievement"})
	}

	return c.JSON(fiber.Map{"message": "achievement deleted"})
}

// ===================== VERIFY / REJECT =====================

func (s *AchievementService) VerifyAchievement(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	if token.Claims.(jwt.MapClaims)["role"] != "Dosen Wali" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	if err := s.Repo.Verify(c.Params("id"), token.Claims.(jwt.MapClaims)["sub"].(string)); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to verify achievement"})
	}

	return c.JSON(fiber.Map{"message": "achievement verified"})
}

func (s *AchievementService) RejectAchievement(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	if token.Claims.(jwt.MapClaims)["role"] != "Dosen Wali" {
		return c.SendStatus(fiber.StatusForbidden)
	}

	var body struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&body); err != nil || body.Note == "" {
		return c.Status(400).JSON(fiber.Map{"message": "rejection note is required"})
	}

	if err := s.Repo.Reject(c.Params("id"), body.Note); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to reject achievement"})
	}

	return c.JSON(fiber.Map{"message": "achievement rejected"})
}
