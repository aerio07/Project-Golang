package repository

import (
	"project_uas/database"
)

// =====================
// PUBLIC FUNCTIONS
// =====================

// Admin → lihat semua achievements
func GetAllAchievements() ([]map[string]interface{}, error) {
	query := `
		SELECT id, student_id, status
		FROM achievement_references
		ORDER BY created_at DESC
	`
	return fetchAchievements(query)
}

// Mahasiswa → lihat achievements sendiri
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

// Dosen Wali → lihat achievements mahasiswa bimbingan
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
// INTERNAL HELPER
// =====================

// helper khusus repository (TIDAK diexport)
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
