package model

// ===== Requests =====

type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin123"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ===== Responses (SRS: status + data) =====

type AuthUser struct {
	ID          string   `json:"id" example:"62d58c69-e469-4acb-98bf-4fd4a881c0d0"`
	Username    string   `json:"username" example:"admin"`
	FullName    string   `json:"fullName" example:"Administrator"`
	Role        string   `json:"role" example:"Admin"`
	Permissions []string `json:"permissions" example:"achievement:create,achievement:read,achievement:update,achievement:delete,achievement:verify,user:manage"`
}

type LoginResponse struct {
	Status string            `json:"status" example:"success"`
	Data   LoginResponseData `json:"data"`
}

type LoginResponseData struct {
	Token        string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string   `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User         AuthUser `json:"user"`
}

type RefreshResponse struct {
	Status string              `json:"status" example:"success"`
	Data   RefreshResponseData `json:"data"`
}

type RefreshResponseData struct {
	Token        string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type ProfileResponse struct {
	Status string              `json:"status" example:"success"`
	Data   ProfileResponseData `json:"data"`
}

type ProfileResponseData struct {
	User AuthUser `json:"user"`
}

type LogoutResponse struct {
	Status string             `json:"status" example:"success"`
	Data   LogoutResponseData `json:"data"`
}

type LogoutResponseData struct {
	Message string `json:"message" example:"logout success"`
}

// ===== Internal helper (DB result) =====

type AuthUserInfo struct {
	ID       string
	Username string
	FullName string
	RoleID   string
	RoleName string
	IsActive bool
}
