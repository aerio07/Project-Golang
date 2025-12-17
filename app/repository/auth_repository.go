package repository

import (
	"database/sql"
	"errors"

	"project_uas/app/model"
)

type AuthRepository interface {
	GetUserByIdentifier(identifier string) (*model.UserLogin, error)
	GetPermissionsByRole(roleID string) ([]string, error)

	// new for profile/refresh
	GetUserByID(userID string) (*model.AuthUserInfo, bool, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

// =====================
// USER
// =====================

func (r *authRepository) GetUserByIdentifier(identifier string) (*model.UserLogin, error) {
	query := `
		SELECT
			u.id,
			u.username,
			u.full_name,
			u.password_hash,
			u.role_id,
			u.is_active,
			r.name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.username = $1 OR u.email = $1
		LIMIT 1
	`

	var user model.UserLogin
	err := r.db.QueryRow(query, identifier).Scan(
		&user.ID,
		&user.Username,
		&user.FullName,
		&user.PasswordHash,
		&user.RoleID,
		&user.IsActive,
		&user.RoleName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) GetUserByID(userID string) (*model.AuthUserInfo, bool, error) {
	query := `
		SELECT
			u.id,
			u.username,
			u.full_name,
			u.role_id,
			u.is_active,
			r.name
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1
		LIMIT 1
	`

	var u model.AuthUserInfo
	err := r.db.QueryRow(query, userID).Scan(
		&u.ID,
		&u.Username,
		&u.FullName,
		&u.RoleID,
		&u.IsActive,
		&u.RoleName,
	)

	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return &u, true, nil
}

// =====================
// PERMISSIONS
// =====================

func (r *authRepository) GetPermissionsByRole(roleID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.name
	`

	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms = append(perms, name)
	}
	return perms, nil
}
