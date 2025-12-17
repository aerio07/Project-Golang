package model

// Request bodies
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// Response bodies (sesuai gaya SRS: status + data)
type AuthUser struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	FullName    string   `json:"fullName"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type LoginResponse struct {
	Status string           `json:"status"`
	Data   LoginResponseData `json:"data"`
}

type LoginResponseData struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refreshToken"`
	User         AuthUser `json:"user"`
}

type RefreshResponse struct {
	Status string            `json:"status"`
	Data   RefreshResponseData `json:"data"`
}

type RefreshResponseData struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type ProfileResponse struct {
	Status string            `json:"status"`
	Data   ProfileResponseData `json:"data"`
}

type ProfileResponseData struct {
	User AuthUser `json:"user"`
}

type LogoutResponse struct {
	Status string           `json:"status"`
	Data   LogoutResponseData `json:"data"`
}

type LogoutResponseData struct {
	Message string `json:"message"`
}

// Internal helper struct from DB
type AuthUserInfo struct {
	ID       string
	Username string
	FullName string
	RoleID   string
	RoleName string
	IsActive bool
}
