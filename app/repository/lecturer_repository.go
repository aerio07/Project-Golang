package repository

import (
	"database/sql"

	"project_uas/app/model"
)

type LecturerRepository interface {
	List(limit, offset int) ([]model.Lecturer, error)
	GetByUserID(userID string) (lecturerID string, ok bool, err error)
	GetAdvisees(lecturerID string, limit, offset int) ([]model.Student, error)
}

type lecturerRepository struct{ db *sql.DB }

func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepository{db: db}
}

func (r *lecturerRepository) List(limit, offset int) ([]model.Lecturer, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.db.Query(`
		SELECT l.id, l.user_id, l.lecturer_id, COALESCE(l.department,''), l.created_at,
		       u.full_name, u.email
		FROM lecturers l
		JOIN users u ON u.id = l.user_id
		ORDER BY l.created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Lecturer
	for rows.Next() {
		var l model.Lecturer
		if err := rows.Scan(&l.ID, &l.UserID, &l.LecturerID, &l.Department, &l.CreatedAt, &l.FullName, &l.Email); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, nil
}

func (r *lecturerRepository) GetByUserID(userID string) (string, bool, error) {
	var id string
	err := r.db.QueryRow(`SELECT id FROM lecturers WHERE user_id=$1`, userID).Scan(&id)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return id, true, nil
}

func (r *lecturerRepository) GetAdvisees(lecturerID string, limit, offset int) ([]model.Student, error) {
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
		WHERE s.advisor_id=$1
		ORDER BY s.created_at DESC
		LIMIT $2 OFFSET $3
	`, lecturerID, limit, offset)
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
