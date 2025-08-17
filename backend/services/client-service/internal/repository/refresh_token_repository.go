package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"microbank/client-service/internal/models"
)

// RefreshTokenRepositoryImpl handles all database operations related to refresh tokens
type RefreshTokenRepositoryImpl struct {
	db *PostgresDB
}

// NewRefreshTokenRepository creates a new refresh token repository
func NewRefreshTokenRepository(db *PostgresDB) RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{db: db}
}

// Create creates a new refresh token in the database
func (r *RefreshTokenRepositoryImpl) Create(refreshToken *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(
		query,
		refreshToken.ID,
		refreshToken.UserID,
		refreshToken.TokenHash,
		refreshToken.ExpiresAt,
		refreshToken.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	return nil
}

// GetByToken retrieves a refresh token by its hash
func (r *RefreshTokenRepositoryImpl) GetByToken(tokenHash string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens WHERE token_hash = $1`

	refreshToken := &models.RefreshToken{}
	err := r.db.QueryRow(query, tokenHash).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return refreshToken, nil
}

// GetByUserID retrieves all refresh tokens for a specific user
func (r *RefreshTokenRepositoryImpl) GetByUserID(userID uuid.UUID) ([]models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query refresh tokens: %w", err)
	}
	defer rows.Close()

	var refreshTokens []models.RefreshToken
	for rows.Next() {
		var refreshToken models.RefreshToken
		err := rows.Scan(
			&refreshToken.ID,
			&refreshToken.UserID,
			&refreshToken.TokenHash,
			&refreshToken.ExpiresAt,
			&refreshToken.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan refresh token row: %w", err)
		}
		refreshTokens = append(refreshTokens, refreshToken)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over refresh token rows: %w", err)
	}

	return refreshTokens, nil
}

// Delete deletes a specific refresh token
func (r *RefreshTokenRepositoryImpl) Delete(id uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh token not found for deletion")
	}

	return nil
}

// DeleteByUserID deletes all refresh tokens for a specific user
func (r *RefreshTokenRepositoryImpl) DeleteByUserID(userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh tokens by user ID: %w", err)
	}

	return nil
}

// DeleteExpired deletes all expired refresh tokens
func (r *RefreshTokenRepositoryImpl) DeleteExpired() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`

	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired refresh tokens: %w", err)
	}

	return nil
}

// CleanupExpiredTokens removes expired tokens (should be called periodically)
func (r *RefreshTokenRepositoryImpl) CleanupExpiredTokens() error {
	return r.DeleteExpired()
}
