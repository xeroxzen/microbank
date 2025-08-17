package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email" binding:"required,email"`
	Name         string    `json:"name" db:"name" binding:"required,min=2,max=100"`
	PasswordHash string    `json:"-" db:"password_hash"`
	IsBlacklisted bool     `json:"is_blacklisted" db:"is_blacklisted"`
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserRegistration represents the data needed to register a new user
type UserRegistration struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserLogin represents the data needed to login a user
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserProfile represents the user profile data that can be updated
type UserProfile struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// UserResponse represents the user data sent in responses (excludes sensitive info)
type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	IsBlacklisted bool     `json:"is_blacklisted"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RefreshToken represents a refresh token for JWT authentication
type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ToResponse converts a User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		IsBlacklisted: u.IsBlacklisted,
		IsAdmin:      u.IsAdmin,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// IsValid checks if the user is valid for operations
func (u *User) IsValid() bool {
	return !u.IsBlacklisted && u.ID != uuid.Nil
}
