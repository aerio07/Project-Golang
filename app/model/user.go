package model

type UserLogin struct {
	ID           string
	Username     string
	FullName     string
	PasswordHash string
	RoleID       string
	RoleName     string
	IsActive     bool
}
