package repository

import (
	"github.com/google/uuid"
	"microbank/client-service/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	UpdateBlacklistStatus(userID uuid.UUID, isBlacklisted bool) error
	GetAllUsers() ([]models.User, error)
	DeleteUser(id uuid.UUID) error
	UserExists(email string) (bool, error)
}

// RefreshTokenRepository defines the interface for refresh token operations
type RefreshTokenRepository interface {
	Create(refreshToken *models.RefreshToken) error
	GetByToken(tokenHash string) (*models.RefreshToken, error)
	GetByUserID(userID uuid.UUID) ([]models.RefreshToken, error)
	Delete(id uuid.UUID) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteExpired() error
	CleanupExpiredTokens() error
}
