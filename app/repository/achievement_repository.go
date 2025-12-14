package repository

import (
	"project_uas/database"

	"github.com/google/uuid"
)

// =====================
// READ FUNCTIONS
// =====================

func GetAllAchievements() ([]map[string]interface{}, error) {
	query := `
		SELECT id, student_id, status
		FROM achievement_references
		ORDER BY created_at DESC
	`
	return fetchAchievements(query)
}

func GetAchievementsByStudent(userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT ar.id, ar.student_id, ar.status
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		WHERE s.user_id = $1
		ORDER BY ar.created_at DESC
	`
	return fetchAchievements(query, userID)
}

func GetAchievementsBySupervisor(userID string) ([]map[string]interface{}, error) {
	query := `
		SELECT ar.id, ar.student_id, ar.status
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		WHERE s.supervisor_user_id = $1
		ORDER BY ar.created_at DESC
	`
	return fetchAchievements(query, userID)
}

// =====================
// CREATE FUNCTION
// =====================

func CreateAchievementReference(studentID string) error {
	query := `
		INSERT INTO achievement_references (
			id, student_id, status
		) VALUES (
			$1, $2, 'submitted'
		)
	`

	_, err := database.DB.Exec(
		query,
		uuid.New().String(),
		studentID,
	)

	return err
}

// =====================
// INTERNAL HELPER
// =====================

func fetchAchievements(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var id, studentID, status string

		if err := rows.Scan(&id, &studentID, &status); err != nil {
			return nil, err
		}

		results = append(results, map[string]interface{}{
			"id":         id,
			"student_id": studentID,
			"status":     status,
		})
	}

	return results, nil
}
