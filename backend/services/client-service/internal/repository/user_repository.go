package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"microbank/client-service/internal/models"
)

// UserRepositoryImpl handles all database operations related to users
type UserRepositoryImpl struct {
	db *PostgresDB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *PostgresDB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

// CreateUser creates a new user in the database
func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, name, password_hash, is_blacklisted, is_admin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.IsBlacklisted,
		user.IsAdmin,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepositoryImpl) GetUserByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, is_blacklisted, is_admin, created_at, updated_at
		FROM users WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.IsBlacklisted,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (r *UserRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, name, password_hash, is_blacklisted, is_admin, created_at, updated_at
		FROM users WHERE email = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.IsBlacklisted,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user's information
func (r *UserRepositoryImpl) UpdateUser(user *models.User) error {
	query := `
		UPDATE users 
		SET name = $1, updated_at = $2
		WHERE id = $3`

	user.UpdatedAt = time.Now()

	result, err := r.db.Exec(query, user.Name, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found for update")
	}

	return nil
}

// UpdateBlacklistStatus updates a user's blacklist status
func (r *UserRepositoryImpl) UpdateBlacklistStatus(userID uuid.UUID, isBlacklisted bool) error {
	query := `
		UPDATE users 
		SET is_blacklisted = $1, updated_at = $2
		WHERE id = $3`

	result, err := r.db.Exec(query, isBlacklisted, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update blacklist status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found for blacklist update")
	}

	return nil
}

// GetAllUsers retrieves all users (for admin purposes)
func (r *UserRepositoryImpl) GetAllUsers() ([]models.User, error) {
	query := `
		SELECT id, email, name, password_hash, is_blacklisted, is_admin, created_at, updated_at
		FROM users
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.PasswordHash,
			&user.IsBlacklisted,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %w", err)
	}

	return users, nil
}

// DeleteUser deletes a user from the database
func (r *UserRepositoryImpl) DeleteUser(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found for deletion")
	}

	return nil
}

// UserExists checks if a user with the given email exists
func (r *UserRepositoryImpl) UserExists(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return exists, nil
}
