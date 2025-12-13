package service

import (
	"project_uas/app/repository"
	"project_uas/database"
	"project_uas/utils"
	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid request",
		})
	}

	// ambil user dari DB
	user, err := repository.GetUserByUsernameOrEmail(database.DB, req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": "invalid username or password",
		})
	}

	// cek aktif
	if !user.IsActive {
		return c.Status(403).JSON(fiber.Map{
			"message": "user is inactive",
		})
	}

	// cek password
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		return c.Status(401).JSON(fiber.Map{
			"message": "invalid username or password",
		})
	}

	// sementara response sukses
	// ambil permissions berdasarkan role
perms, err := repository.GetPermissionsByRoleID(database.DB, user.RoleID)
if err != nil {
	return c.Status(500).JSON(fiber.Map{
		"message": "failed to load permissions",
	})
}

// generate JWT
token, err := utils.GenerateToken(user.ID, user.RoleName, perms)
if err != nil {
	return c.Status(500).JSON(fiber.Map{
		"message": err.Error(),
	})
}

// response FINAL (LOGIN + JWT)
return c.JSON(fiber.Map{
	"message":     "login success",
	"user_id":     user.ID,
	"role":        user.RoleName,
	"token":       token,
	"permissions": perms,
})


}
