package service

import (
	"project_uas/app/repository"
	"project_uas/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthService struct {
	Repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// =====================
// LOGIN
// =====================

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	user, err := s.Repo.GetUserByIdentifier(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "invalid username or password",
		})
	}

	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{
			"message": "user is inactive",
		})
	}

	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return c.Status(401).JSON(fiber.Map{
			"message": "invalid username or password",
		})
	}

	perms, err := s.Repo.GetPermissionsByRole(user.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "failed to load permissions",
		})
	}

	token, err := utils.GenerateToken(user.ID, user.RoleName, perms)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "login success",
		"user_id":     user.ID,
		"role":        user.RoleName,
		"token":       token,
		"permissions": perms,
	})
}
