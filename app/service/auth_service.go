package service

import (
	"strings"

	"project_uas/app/model"
	"project_uas/app/repository"
	"project_uas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	Repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

// =====================
// LOGIN (SRS)
// =====================

func (s *AuthService) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "invalid request"})
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "username and password are required"})
	}

	user, err := s.Repo.GetUserByIdentifier(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "invalid username or password"})
	}
	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "user is inactive"})
	}
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "invalid username or password"})
	}

	perms, err := s.Repo.GetPermissionsByRole(user.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to load permissions"})
	}

	accessToken, err := utils.GenerateToken(user.ID, user.RoleName, perms)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to generate refresh token"})
	}

	resp := model.LoginResponse{
		Status: "success",
		Data: model.LoginResponseData{
			Token:        accessToken,
			RefreshToken: refreshToken,
			User: model.AuthUser{
				ID:          user.ID,
				Username:    user.Username,
				FullName:    user.FullName,
				Role:        user.RoleName,
				Permissions: perms,
			},
		},
	}
	return c.JSON(resp)
}

// =====================
// REFRESH (stateless)
// =====================

func (s *AuthService) Refresh(c *fiber.Ctx) error {
	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "invalid request"})
	}

	req.RefreshToken = strings.TrimSpace(req.RefreshToken)
	if req.RefreshToken == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "refreshToken is required"})
	}

	claims, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "invalid refresh token"})
	}

	userID := claims.Subject
	info, ok, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to load user"})
	}
	if !ok {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "user not found"})
	}
	if !info.IsActive {
		return c.Status(403).JSON(fiber.Map{"status": "error", "message": "user is inactive"})
	}

	perms, err := s.Repo.GetPermissionsByRole(info.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to load permissions"})
	}

	newAccess, err := utils.GenerateToken(info.ID, info.RoleName, perms)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to generate token"})
	}

	newRefresh, err := utils.GenerateRefreshToken(info.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to generate refresh token"})
	}

	return c.JSON(model.RefreshResponse{
		Status: "success",
		Data: model.RefreshResponseData{
			Token:        newAccess,
			RefreshToken: newRefresh,
		},
	})
}

// =====================
// LOGOUT (stateless)
// =====================

func (s *AuthService) Logout(c *fiber.Ctx) error {
	// Karena stateless: server tidak revoke token.
	// Client WAJIB hapus token yang tersimpan.
	return c.JSON(model.LogoutResponse{
		Status: "success",
		Data:   model.LogoutResponseData{Message: "logout success"},
	})
}

// =====================
// PROFILE (pakai access token)
// =====================

func (s *AuthService) Profile(c *fiber.Ctx) error {
	tokAny := c.Locals("user")
	token, ok := tokAny.(*jwt.Token)
	if !ok || token == nil {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}

	userID, _ := claims["sub"].(string)
	if strings.TrimSpace(userID) == "" {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "unauthorized"})
	}

	info, ok, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to load user"})
	}
	if !ok {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "user not found"})
	}

	perms, err := s.Repo.GetPermissionsByRole(info.RoleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "failed to load permissions"})
	}

	return c.JSON(model.ProfileResponse{
		Status: "success",
		Data: model.ProfileResponseData{
			User: model.AuthUser{
				ID:          info.ID,
				Username:    info.Username,
				FullName:    info.FullName,
				Role:        info.RoleName,
				Permissions: perms,
			},
		},
	})
}
