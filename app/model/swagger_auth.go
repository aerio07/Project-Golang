package model

// ===== AUTH Swagger Helper Types =====

// Karena di Auth kamu sering balikin:
// {"status":"error","message":"..."}
// maka bikin schema khusus biar Swagger akurat.
type AuthErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"invalid username or password"`
}

// Optional: buat response unauthorized standar
type AuthUnauthorizedResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"unauthorized"`
}
