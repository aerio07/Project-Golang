package model

// ===== Swagger Responses untuk Lecturer Module =====
// NOTE: Student & Lecturer sudah ada di model kamu, jadi DI SINI jangan buat ulang type-nya.

type LecturerListResponse struct {
	Data []Lecturer `json:"data"`
}

type AdviseeListResponse struct {
	Data []Student `json:"data"`
}
