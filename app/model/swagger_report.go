package model

// Wrapper response: {"data": ...}
type ReportStatisticsResponse struct {
	Data *AchievementStatistics `json:"data"`
}

type ReportStudentResponse struct {
	Data *AchievementStatistics `json:"data"`
}
