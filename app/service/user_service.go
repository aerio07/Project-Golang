package service

import (
	"strings"


	"project_uas/app/repository"
	"project_uas/utils"

	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// GetUsers godoc
// @Summary List users
// @Description List user (admin management). Bisa search via q
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param q query string false "Search by username/email/full_name"
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} model.UserListResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users [get]
func (s *UserService) GetUsers(c *fiber.Ctx) error {
	q := c.Query("q", "")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	users, err := s.Repo.List(q, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch users"})
	}
	return c.JSON(fiber.Map{"data": users})
}

// GetUser godoc
// @Summary Get user detail
// @Description Detail user by id
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID (uuid)"
// @Success 200 {object} model.UserDetailResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/{id} [get]
func (s *UserService) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	u, ok, err := s.Repo.GetByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch user"})
	}
	if !ok {
		return c.Status(404).JSON(fiber.Map{"message": "user not found"})
	}
	return c.JSON(fiber.Map{"data": u})
}

type createUserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	RoleName string `json:"roleName"` // "Admin" | "Mahasiswa" | "Dosen Wali"
}

// CreateUser godoc
// @Summary Create user
// @Description Buat user baru + set role
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.UserCreateRequest true "Request body"
// @Success 201 {object} model.UserCreateResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users [post]
func (s *UserService) CreateUser(c *fiber.Ctx) error {
	var body createUserReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	body.Username = strings.TrimSpace(body.Username)
	body.Email = strings.TrimSpace(body.Email)
	body.Password = strings.TrimSpace(body.Password)
	body.FullName = strings.TrimSpace(body.FullName)
	body.RoleName = strings.TrimSpace(body.RoleName)

	if body.Username == "" || body.Email == "" || body.Password == "" || body.FullName == "" || body.RoleName == "" {
		return c.Status(400).JSON(fiber.Map{"message": "username, email, password, full_name, roleName are required"})
	}

	hash, err := utils.HashPassword(body.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to hash password"})
	}

	userID, _, err := s.Repo.CreateUserWithRole(body.Username, body.Email, hash, body.FullName, body.RoleName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to create user"})
	}

	return c.Status(201).JSON(fiber.Map{"data": fiber.Map{"id": userID}})
}

type updateUserReq struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	FullName *string `json:"full_name"`
	IsActive *bool   `json:"is_active"`
}

// UpdateUser godoc
// @Summary Update user
// @Description Update sebagian field (partial update) via COALESCE
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID (uuid)"
// @Param body body model.UserUpdateRequest true "Request body"
// @Success 200 {object} model.MessageResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/{id} [put]
func (s *UserService) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	var body updateUserReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "invalid body"})
	}

	if err := s.Repo.UpdateUser(id, body.Username, body.Email, body.FullName, body.IsActive); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to update user"})
	}
	return c.JSON(fiber.Map{"message": "user updated"})
}

// DeleteUser godoc
// @Summary Deactivate user
// @Description Soft delete user (set is_active=false)
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID (uuid)"
// @Success 200 {object} model.MessageResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/{id} [delete]
func (s *UserService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.Repo.Deactivate(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to deactivate user"})
	}
	return c.JSON(fiber.Map{"message": "user deactivated"})
}

type assignRoleReq struct {
	RoleName string `json:"roleName"`
}

// AssignRole godoc
// @Summary Assign role
// @Description Ganti role user berdasarkan roleName
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID (uuid)"
// @Param body body model.UserAssignRoleRequest true "Request body"
// @Success 200 {object} model.MessageResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.AuthUnauthorizedResponse
// @Failure 403 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/{id}/role [put]
func (s *UserService) AssignRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var body assignRoleReq
	if err := c.BodyParser(&body); err != nil || strings.TrimSpace(body.RoleName) == "" {
		return c.Status(400).JSON(fiber.Map{"message": "roleName is required"})
	}

	if err := s.Repo.AssignRole(id, strings.TrimSpace(body.RoleName)); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to assign role"})
	}
	return c.JSON(fiber.Map{"message": "role updated"})
}
