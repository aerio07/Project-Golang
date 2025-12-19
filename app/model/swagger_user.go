package model

// ===== Users =====

type UserListResponse struct {
	Data []User `json:"data"`
}

type UserDetailResponse struct {
	Data User `json:"data"`
}

type UserCreateRequest struct {
	Username string `json:"username" example:"admin1"`
	Email    string `json:"email" example:"admin1@mail.com"`
	Password string `json:"password" example:"admin123"`
	FullName string `json:"full_name" example:"Admin Satu"`
	RoleName string `json:"roleName" example:"Admin"` // Admin | Mahasiswa | Dosen Wali
}

type UserCreateResponse struct {
	Data struct {
		ID string `json:"id" example:"62d58c69-e469-4acb-98bf-4fd4a881c0d0"`
	} `json:"data"`
}

type UserUpdateRequest struct {
	Username *string `json:"username,omitempty" example:"admin2"`
	Email    *string `json:"email,omitempty" example:"admin2@mail.com"`
	FullName *string `json:"full_name,omitempty" example:"Admin Dua"`
	IsActive *bool   `json:"is_active,omitempty" example:"true"`
}

type UserAssignRoleRequest struct {
	RoleName string `json:"roleName" example:"Mahasiswa"`
}
