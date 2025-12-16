package model

type CountByKey struct {
	Key   string `json:"key" bson:"_id"`
	Total int64  `json:"total" bson:"total"`
}

type TopStudent struct {
	StudentID string `json:"student_id" bson:"_id"`
	Points    int64  `json:"points" bson:"points"`
	Total     int64  `json:"total" bson:"total"`
}

type AchievementStatistics struct {
	TotalPerType    []CountByKey `json:"total_per_type"`
	TotalPerPeriod  []CountByKey `json:"total_per_period"`
	TopStudents     []TopStudent `json:"top_students"`
	LevelDist       []CountByKey `json:"level_distribution"`
	Scope           string       `json:"scope"` // "all" / "student" / "advisees"
	FilteredStudents int         `json:"filtered_students"`
}
