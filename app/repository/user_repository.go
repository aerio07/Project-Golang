package repository

import (
	"database/sql"
	"errors"
	"time"

	"project_uas/app/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	List(search string, limit, offset int) ([]model.User, error)
	GetByID(id string) (*model.User, bool, error)

	CreateUserWithRole(
		username, email, passwordHash, fullName, roleName string,
	) (userID string, roleID string, err error)

	UpdateUser(id string, username, email, fullName *string, isActive *bool) error

	Deactivate(id string) error
	
	GetRoleIDByName(roleName string) (string, bool, error)
	AssignRole(userID, roleName string) error
}

type userRepository struct{ db *sql.DB }

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) List(search string, limit, offset int) ([]model.User, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	q := `
		SELECT u.id, u.username, u.email, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE ($1 = '' OR u.username ILIKE '%' || $1 || '%' OR u.email ILIKE '%' || $1 || '%' OR u.full_name ILIKE '%' || $1 || '%')
		ORDER BY u.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(q, search, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func (r *userRepository) GetByID(id string) (*model.User, bool, error) {
	q := `
		SELECT u.id, u.username, u.email, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id=$1
	`
	var u model.User
	err := r.db.QueryRow(q, id).Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &u, true, nil
}

func (r *userRepository) GetRoleIDByName(roleName string) (string, bool, error) {
	var id string
	err := r.db.QueryRow(`SELECT id FROM roles WHERE name=$1`, roleName).Scan(&id)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return id, true, nil
}

func (r *userRepository) CreateUserWithRole(username, email, passwordHash, fullName, roleName string) (string, string, error) {
	roleID, ok, err := r.GetRoleIDByName(roleName)
	if err != nil {
		return "", "", err
	}
	if !ok {
		return "", "", errors.New("role not found")
	}

	var userID string
	now := time.Now()

	err = r.db.QueryRow(`
		INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, true, $6, $6)
		RETURNING id
	`, username, email, passwordHash, fullName, roleID, now).Scan(&userID)

	return userID, roleID, err
}

func (r *userRepository) UpdateUser(id string, username, email, fullName *string, isActive *bool) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET username   = COALESCE($2, username),
		    email      = COALESCE($3, email),
		    full_name  = COALESCE($4, full_name),
		    is_active  = COALESCE($5, is_active),
		    updated_at = NOW()
		WHERE id = $1
	`, id, username, email, fullName, isActive)

	return err
}


func (r *userRepository) Deactivate(id string) error {
	_, err := r.db.Exec(`UPDATE users SET is_active=false, updated_at=NOW() WHERE id=$1`, id)
	return err
}

func (r *userRepository) AssignRole(userID, roleName string) error {
	roleID, ok, err := r.GetRoleIDByName(roleName)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("role not found")
	}

	_, err = r.db.Exec(`UPDATE users SET role_id=$2, updated_at=NOW() WHERE id=$1`, userID, roleID)
	return err
}
