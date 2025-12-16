package repository

import (
	"database/sql"
	"errors"
	"time"

	"project_uas/app/model"

	"github.com/google/uuid"
)

var (
	ErrNotFoundOrForbidden = errors.New("not found or forbidden")
)

type AchievementRepository interface {
	// list
	GetAll() ([]model.Achievement, error)
	GetByStudent(userID string) ([]model.Achievement, error)
	GetBySupervisor(userID string) ([]model.Achievement, error)

	// helpers
	GetStudentIDByUserID(userID string) (string, bool, error)
	GetStatusByID(refID string) (string, bool, error)

	// create
	CreateDraft(studentID string) (string, error)
	SetMongoID(refID string, mongoID string) error

	// detail access
	GetRefForDetailStudent(refID, userID string) (*string, string, bool, error)
	GetRefForDetailSupervisor(refID, userID string) (*string, string, bool, error)
	GetRefForDetailAdmin(refID string) (*string, string, bool, error)

	// actions
	Submit(id, userID string) error
	CanDelete(id, userID string) (bool, error)
	SoftDelete(id string) error
	Verify(id, verifierID string) error
	Reject(id, note, rejecterID string) error

	// history (implicit)
	GetImplicitHistory(refID string) ([]model.AchievementHistory, error)
}

type achievementRepository struct {
	db *sql.DB
}

func NewAchievementRepository(db *sql.DB) AchievementRepository {
	return &achievementRepository{db: db}
}

//
// ===== LIST =====
//

func (r *achievementRepository) GetAll() ([]model.Achievement, error) {
	return r.fetch(`
		SELECT id, student_id, status, created_at
		FROM achievement_references
		WHERE status != 'deleted'
		ORDER BY created_at DESC
	`)
}

func (r *achievementRepository) GetByStudent(userID string) ([]model.Achievement, error) {
	return r.fetch(`
		SELECT ar.id, ar.student_id, ar.status, ar.created_at
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		WHERE s.user_id = $1 AND ar.status != 'deleted'
		ORDER BY ar.created_at DESC
	`, userID)
}

func (r *achievementRepository) GetBySupervisor(userID string) ([]model.Achievement, error) {
	return r.fetch(`
		SELECT ar.id, ar.student_id, ar.status, ar.created_at
		FROM achievement_references ar
		JOIN students s ON s.id = ar.student_id
		WHERE s.advisor_id = $1 AND ar.status != 'deleted'
		ORDER BY ar.created_at DESC
	`, userID)
}

func (r *achievementRepository) fetch(query string, args ...any) ([]model.Achievement, error) {
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

//
// ===== HELPERS =====
//

func (r *achievementRepository) GetStudentIDByUserID(userID string) (string, bool, error) {
	var id string
	err := r.db.QueryRow(`SELECT id FROM students WHERE user_id=$1`, userID).Scan(&id)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	return id, err == nil, err
}

func (r *achievementRepository) GetStatusByID(refID string) (string, bool, error) {
	var status string
	err := r.db.QueryRow(`
		SELECT status FROM achievement_references WHERE id=$1
	`, refID).Scan(&status)

	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return status, true, nil
}


//
// ===== CREATE =====
//

func (r *achievementRepository) CreateDraft(studentID string) (string, error) {
	id := uuid.New().String()
	_, err := r.db.Exec(`
		INSERT INTO achievement_references (id, student_id, status, created_at)
		VALUES ($1, $2, 'draft', NOW())
	`, id, studentID)
	return id, err
}

func (r *achievementRepository) SetMongoID(refID string, mongoID string) error {
	_, err := r.db.Exec(`
		UPDATE achievement_references SET mongo_achievement_id=$2 WHERE id=$1
	`, refID, mongoID)
	return err
}

//
// ===== DETAIL =====
//

func (r *achievementRepository) GetRefForDetailStudent(refID, userID string) (*string, string, bool, error) {
	return r.getDetail(`
		JOIN students s ON s.id = ar.student_id
		WHERE ar.id=$1 AND s.user_id=$2 AND ar.status!='deleted'
	`, refID, userID)
}

func (r *achievementRepository) GetRefForDetailSupervisor(refID, userID string) (*string, string, bool, error) {
	return r.getDetail(`
		JOIN students s ON s.id = ar.student_id
		WHERE ar.id=$1 AND s.advisor_id=$2 AND ar.status!='deleted'
	`, refID, userID)
}

func (r *achievementRepository) GetRefForDetailAdmin(refID string) (*string, string, bool, error) {
	return r.getDetail(`WHERE ar.id=$1 AND ar.status!='deleted'`, refID)
}

func (r *achievementRepository) getDetail(where string, args ...any) (*string, string, bool, error) {
	query := `
		SELECT ar.mongo_achievement_id, ar.status
		FROM achievement_references ar
	` + where

	var mid sql.NullString
	var status string

	err := r.db.QueryRow(query, args...).Scan(&mid, &status)
	if err == sql.ErrNoRows {
		return nil, "", false, nil
	}
	if err != nil {
		return nil, "", false, err
	}
	if mid.Valid {
		return &mid.String, status, true, nil
	}
	return nil, status, true, nil
}

//
// ===== ACTIONS =====
//

func (r *achievementRepository) Submit(id, userID string) error {
	res, err := r.db.Exec(`
		UPDATE achievement_references ar
		SET status='submitted', submitted_at=NOW()
		FROM students s
		WHERE ar.id=$1 AND ar.student_id=s.id
		  AND s.user_id=$2 AND ar.status='draft'
	`, id, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFoundOrForbidden
	}
	return nil
}

func (r *achievementRepository) CanDelete(id, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM achievement_references ar
		JOIN students s ON s.id=ar.student_id
		WHERE ar.id=$1 AND s.user_id=$2 AND ar.status='draft'
	`, id, userID).Scan(&count)
	return count > 0, err
}

func (r *achievementRepository) SoftDelete(id string) error {
	res, err := r.db.Exec(`
		UPDATE achievement_references
		SET status='deleted'
		WHERE id=$1 AND status='draft'
	`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFoundOrForbidden
	}
	return nil
}

func (r *achievementRepository) Verify(id, verifierUserID string) error {
	res, err := r.db.Exec(`
		UPDATE achievement_references ar
		SET status='verified',
		    verified_at=NOW(),
		    verified_by=$2
		FROM students s
		JOIN lecturers l ON l.id = s.advisor_id
		WHERE ar.id = $1
		  AND ar.student_id = s.id
		  AND l.user_id = $2          -- ✅ BENAR
		  AND ar.status = 'submitted'
	`, id, verifierUserID)

	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFoundOrForbidden
	}
	return nil
}


func (r *achievementRepository) Reject(id, note, rejecterUserID string) error {
	res, err := r.db.Exec(`
		UPDATE achievement_references ar
		SET status='rejected',
		    rejection_note=$3
		FROM students s
		JOIN lecturers l ON l.id = s.advisor_id
		WHERE ar.id = $1
		  AND ar.student_id = s.id
		  AND l.user_id = $2          -- ✅ BENAR
		  AND ar.status = 'submitted'
	`, id, rejecterUserID, note)

	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFoundOrForbidden
	}
	return nil
}


//
// ===== HISTORY (IMPLICIT) =====
//

func (r *achievementRepository) GetImplicitHistory(refID string) ([]model.AchievementHistory, error) {
	var (
		status      string
		createdAt   time.Time
		submittedAt sql.NullTime
		verifiedAt  sql.NullTime
	)

	err := r.db.QueryRow(`
		SELECT status, created_at, submitted_at, verified_at
		FROM achievement_references
		WHERE id = $1
	`, refID).Scan(
		&status,
		&createdAt,
		&submittedAt,
		&verifiedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFoundOrForbidden
	}
	if err != nil {
		return nil, err
	}

	history := []model.AchievementHistory{
		{
			Status: "draft",
			At:     createdAt,
		},
	}

	if submittedAt.Valid {
		history = append(history, model.AchievementHistory{
			Status: "submitted",
			At:     submittedAt.Time,
		})
	}

	if status == "verified" && verifiedAt.Valid {
		history = append(history, model.AchievementHistory{
			Status: "verified",
			At:     verifiedAt.Time,
		})
	}

	if status == "rejected" {
		at := createdAt
		if submittedAt.Valid {
			at = submittedAt.Time
		}

		history = append(history, model.AchievementHistory{
			Status: "rejected",
			At:     at,
		})
	}

	return history, nil
}


