package repository

import (
	"database/sql"
	"errors"

	"project_uas/app/model"
)

func GetUserByUsernameOrEmail(db *sql.DB, identifier string) (*model.UserLogin, error) {
	query := `
		SELECT
			u.id,
			u.username,
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

	err := db.QueryRow(query, identifier).Scan(
		&user.ID,
		&user.Username,
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
