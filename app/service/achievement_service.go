package service

import (
	"project_uas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type CreateAchievementRequest struct {
	StudentID string `json:"student_id"`
}

type RejectAchievementRequest struct {
	Note string `json:"note"`
}

// =====================
// GET ACHIEVEMENTS
// =====================

func GetAchievements(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userID := claims["sub"].(string)
	role := claims["role"].(string)

	var data interface{}
	var err error

	switch role {
	case "Admin":
		data, err = repository.GetAllAchievements()
	case "Mahasiswa":
		data, err = repository.GetAchievementsByStudent(userID)
	case "Dosen Wali":
		data, err = repository.GetAchievementsBySupervisor(userID)
	default:
		return c.Status(403).JSON(fiber.Map{
			"message": "role not allowed",
		})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to fetch achievements",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

// =====================
// CREATE ACHIEVEMENT
// =====================

func CreateAchievement(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	role := claims["role"].(string)
	if role != "Mahasiswa" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only mahasiswa can submit achievement",
		})
	}

	var req CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid request body",
		})
	}

	if req.StudentID == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "student_id is required",
		})
	}

	if err := repository.CreateAchievementReference(req.StudentID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to submit achievement",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "achievement submitted",
	})
}

// =====================
// VERIFY ACHIEVEMENT
// =====================

func VerifyAchievement(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	role := claims["role"].(string)
	if role != "Dosen Wali" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only dosen wali can verify",
		})
	}

	achievementID := c.Params("id")

	if err := repository.VerifyAchievement(achievementID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to verify achievement",
		})
	}

	return c.JSON(fiber.Map{
		"message": "achievement verified",
	})
}

// =====================
// REJECT ACHIEVEMENT
// =====================

func RejectAchievement(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	role := claims["role"].(string)
	if role != "Dosen Wali" {
		return c.Status(403).JSON(fiber.Map{
			"message": "only dosen wali can reject",
		})
	}

	achievementID := c.Params("id")

	var req RejectAchievementRequest
	if err := c.BodyParser(&req); err != nil || req.Note == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "rejection note is required",
		})
	}

	if err := repository.RejectAchievement(achievementID, req.Note); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to reject achievement",
		})
	}

	return c.JSON(fiber.Map{
		"message": "achievement rejected",
	})
}
