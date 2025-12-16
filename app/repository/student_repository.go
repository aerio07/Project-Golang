package repository

import (
	"database/sql"

	"project_uas/app/model"
)

type StudentRepository interface {
	List(limit, offset int) ([]model.Student, error)
	GetByID(id string) (*model.Student, bool, error)
	GetByUserID(userID string) (studentID string, ok bool, err error)

	SetAdvisor(studentID, advisorLecturerID string) error
	GetAchievements(studentID string, limit, offset int) ([]model.Achievement, error)

	// access checks
	IsAdvisee(studentID, lecturerUserID string) (bool, error)
}

type studentRepository struct{ db *sql.DB }

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) List(limit, offset int) ([]model.Student, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.db.Query(`
		SELECT s.id, s.user_id, s.student_id, s.program_study, s.academic_year, s.advisor_id, s.created_at,
		       u.full_name, u.email
		FROM students s
		JOIN users u ON u.id = s.user_id
		ORDER BY s.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Student
	for rows.Next() {
		var s model.Student
		var advisor sql.NullString
		if err := rows.Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &advisor, &s.CreatedAt, &s.FullName, &s.Email); err != nil {
			return nil, err
		}
		if advisor.Valid {
			s.AdvisorID = &advisor.String
		}
		out = append(out, s)
	}
	return out, nil
}

func (r *studentRepository) GetByID(id string) (*model.Student, bool, error) {
	var s model.Student
	var advisor sql.NullString
	err := r.db.QueryRow(`
		SELECT s.id, s.user_id, s.student_id, s.program_study, s.academic_year, s.advisor_id, s.created_at,
		       u.full_name, u.email
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.id=$1
	`, id).Scan(&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &advisor, &s.CreatedAt, &s.FullName, &s.Email)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if advisor.Valid {
		s.AdvisorID = &advisor.String
	}
	return &s, true, nil
}

func (r *studentRepository) GetByUserID(userID string) (string, bool, error) {
	var id string
	err := r.db.QueryRow(`SELECT id FROM students WHERE user_id=$1`, userID).Scan(&id)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return id, true, nil
}

func (r *studentRepository) SetAdvisor(studentID, advisorLecturerID string) error {
	_, err := r.db.Exec(`UPDATE students SET advisor_id=$2 WHERE id=$1`, studentID, advisorLecturerID)
	return err
}

func (r *studentRepository) GetAchievements(studentID string, limit, offset int) ([]model.Achievement, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.db.Query(`
		SELECT id, student_id, status, created_at
		FROM achievement_references
		WHERE student_id=$1 AND status != 'deleted'
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, studentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Achievement
	for rows.Next() {
		var a model.Achievement
		if err := rows.Scan(&a.ID, &a.StudentID, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

func (r *studentRepository) IsAdvisee(studentID, lecturerUserID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*)
		FROM students s
		JOIN lecturers l ON l.id = s.advisor_id
		WHERE s.id=$1 AND l.user_id=$2
	`, studentID, lecturerUserID).Scan(&count)
	return count > 0, err
}
