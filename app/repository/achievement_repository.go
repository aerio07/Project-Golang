package repository

import (
	"database/sql"
	"project_uas/app/model"

	"github.com/google/uuid"
)

type AchievementRepository interface {
	GetAll() ([]model.Achievement, error)
	GetByStudent(userID string) ([]model.Achievement, error)
	GetBySupervisor(userID string) ([]model.Achievement, error)

	CreateDraft(studentID string) error
	Submit(id string, userID string) error

	CanDelete(id string, userID string) (bool, error)
	SoftDelete(id string) error

	Verify(id string, verifierID string) error
	Reject(id string, note string) error
}

type achievementRepository struct {
	db *sql.DB
}

func NewAchievementRepository(db *sql.DB) AchievementRepository {
	return &achievementRepository{db: db}
}

// ===================== READ =====================

func (r *achievementRepository) GetAll() ([]model.Achievement, error) {
	query := `
		SELECT id, student_id, status, created_at
		FROM achievement_references
		WHERE isdelete = false
		ORDER BY created_at DESC
	`
	return r.fetch(query)
}

func (r *achievementRepository) GetByStudent(userID string) ([]model.Achievement, error) {
	query := `
		SELECT ar.id, ar.student_id, ar.status, ar.created_at
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		WHERE s.user_id = $1 AND ar.isdelete = false
		ORDER BY ar.created_at DESC
	`
	return r.fetch(query, userID)
}

func (r *achievementRepository) GetBySupervisor(userID string) ([]model.Achievement, error) {
	query := `
		SELECT ar.id, ar.student_id, ar.status, ar.created_at
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		WHERE s.advisor_id = $1 AND ar.isdelete = false
		ORDER BY ar.created_at DESC
	`
	return r.fetch(query, userID)
}

// ===================== CREATE / SUBMIT =====================

func (r *achievementRepository) CreateDraft(studentID string) error {
	_, err := r.db.Exec(`
		INSERT INTO achievement_references (id, student_id, status, isdelete, created_at)
		VALUES ($1, $2, 'draft', false, NOW())
	`, uuid.New().String(), studentID)
	return err
}

func (r *achievementRepository) Submit(id string, userID string) error {
	_, err := r.db.Exec(`
		UPDATE achievement_references ar
		SET status='submitted', submitted_at=NOW()
		FROM students s
		WHERE ar.id=$1
		  AND ar.student_id=s.id
		  AND s.user_id=$2
		  AND ar.status='draft'
		  AND ar.isdelete=false
	`, id, userID)
	return err
}

// ===================== DELETE =====================

func (r *achievementRepository) CanDelete(id string, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		WHERE ar.id=$1
		  AND s.user_id=$2
		  AND ar.status='draft'
		  AND ar.isdelete=false
	`, id, userID).Scan(&count)
	return count > 0, err
}

func (r *achievementRepository) SoftDelete(id string) error {
	_, err := r.db.Exec(`UPDATE achievement_references SET isdelete=true WHERE id=$1`, id)
	return err
}

// ===================== VERIFY / REJECT =====================

func (r *achievementRepository) Verify(id string, verifierID string) error {
	_, err := r.db.Exec(`
		UPDATE achievement_references
		SET status='verified', verified_at=NOW(), verified_by=$2
		WHERE id=$1 AND status='submitted' AND isdelete=false
	`, id, verifierID)
	return err
}

func (r *achievementRepository) Reject(id string, note string) error {
	_, err := r.db.Exec(`
		UPDATE achievement_references
		SET status='rejected', rejection_note=$2
		WHERE id=$1 AND status='submitted' AND isdelete=false
	`, id, note)
	return err
}

// ===================== HELPER =====================

func (r *achievementRepository) fetch(query string, args ...interface{}) ([]model.Achievement, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Achievement
	for rows.Next() {
		var a model.Achievement
		if err := rows.Scan(&a.ID, &a.StudentID, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

// ambil mongo_achievement_id + status + ownership
func (r *achievementRepository) GetRefForDetail(refID, userID string) (mongoID *string, status string, ok bool, err error) {
	query := `
		SELECT ar.mongo_achievement_id, ar.status
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		WHERE ar.id = $1 AND s.user_id = $2 AND ar.isdelete = false
	`
	var mid sql.NullString
	err = r.db.QueryRow(query, refID, userID).Scan(&mid, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", false, nil
		}
		return nil, "", false, err
	}
	if mid.Valid {
		return &mid.String, status, true, nil
	}
	return nil, status, true, nil
}

func (r *achievementRepository) SetMongoID(refID string, mongoID string) error {
	_, err := r.db.Exec(
		`UPDATE achievement_references SET mongo_achievement_id = $2 WHERE id = $1`,
		refID, mongoID,
	)
	return err
}
