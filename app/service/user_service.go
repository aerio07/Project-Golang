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

func (s *UserService) GetUsers(c *fiber.Ctx) error {
	// optional query: ?q=&limit=&offset=
	q := c.Query("q", "")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	users, err := s.Repo.List(q, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to fetch users"})
	}
	return c.JSON(fiber.Map{"data": users})
}

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

func (s *UserService) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	// aman: soft delete via is_active=false
	if err := s.Repo.Deactivate(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "failed to deactivate user"})
	}
	return c.JSON(fiber.Map{"message": "user deactivated"})
}

type assignRoleReq struct {
	RoleName string `json:"roleName"`
}

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
