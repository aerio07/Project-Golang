package model

type StudentListResponse struct {
	Data []Student `json:"data"`
}

type StudentDetailResponse struct {
	Data Student `json:"data"`
}

type StudentAchievementListResponse struct {
	Data []Achievement `json:"data"`
}

type StudentSetAdvisorRequest struct {
	AdvisorID string `json:"advisor_id" example:"b5a2c9f6-1e7b-4b6c-8b0a-6c9c7d2b0f11"`
}
