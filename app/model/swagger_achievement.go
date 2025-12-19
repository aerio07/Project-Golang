package model

// ===== Common =====

// standar error yang sering kamu pakai: {"message": "..."}
type ErrorResponse struct {
	Message string `json:"message" example:"failed to fetch achievements"`
}

type MessageResponse struct {
	Message string `json:"message" example:"success"`
}

// ===== Requests (Achievements) =====

// AchievementUpsertRequest payload create/update achievement
type AchievementUpsertRequest struct {
	AchievementType string                 `json:"achievementType" example:"competition"`
	Title           string                 `json:"title" example:"Juara 1 Lomba UI/UX"`
	Description     string                 `json:"description" example:"Menang lomba tingkat nasional"`
	Details         map[string]interface{} `json:"details"`
	Tags            []string               `json:"tags" example:"[\"UI\",\"UX\"]"`
	Points          int                    `json:"points" example:"50"`
}

// AchievementRejectRequest payload reject
type AchievementRejectRequest struct {
	Note string `json:"note" example:"Bukti kurang jelas"`
}

// ===== Responses (Achievements) =====

type AchievementCreateResponse struct {
	Data struct {
		ID     string          `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
		Status string          `json:"status" example:"draft"`
		Detail *AchievementMongo `json:"detail"`
	} `json:"data"`
}

type AchievementListResponse struct {
	Data []Achievement `json:"data"`
}

type AchievementDetailResponse struct {
	Data struct {
		ID     string           `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
		Status string           `json:"status" example:"draft"`
		Detail *AchievementMongo `json:"detail"`
	} `json:"data"`
}

type AchievementHistoryResponse struct {
	Data struct {
		CurrentStatus string               `json:"currentStatus" example:"submitted"`
		History       []AchievementHistory `json:"history"`
	} `json:"data"`
}

type AttachmentResponse struct {
	Data Attachment `json:"data"`
}
