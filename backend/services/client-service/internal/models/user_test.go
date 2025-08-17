package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUser_ToResponse(t *testing.T) {
	// Create a test user
	user := &User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "hashedpassword",
		IsBlacklisted: false,
		IsAdmin:      false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Convert to response
	response := user.ToResponse()

	// Verify fields are correctly mapped
	if response.ID != user.ID {
		t.Errorf("Expected ID %v, got %v", user.ID, response.ID)
	}

	if response.Email != user.Email {
		t.Errorf("Expected Email %s, got %s", user.Email, response.Email)
	}

	if response.Name != user.Name {
		t.Errorf("Expected Name %s, got %s", user.Name, response.Name)
	}

	if response.IsBlacklisted != user.IsBlacklisted {
		t.Errorf("Expected IsBlacklisted %v, got %v", user.IsBlacklisted, response.IsBlacklisted)
	}

	if response.IsAdmin != user.IsAdmin {
		t.Errorf("Expected IsAdmin %v, got %v", user.IsAdmin, response.IsAdmin)
	}

	if response.CreatedAt != user.CreatedAt {
		t.Errorf("Expected CreatedAt %v, got %v", user.CreatedAt, response.CreatedAt)
	}

	if response.UpdatedAt != user.UpdatedAt {
		t.Errorf("Expected UpdatedAt %v, got %v", user.UpdatedAt, response.UpdatedAt)
	}

	// Verify password hash is not exposed (UserResponse doesn't have this field)
	// This is correct for security - password hashes should never be exposed
}

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		expected bool
	}{
		{
			name: "valid user",
			user: &User{
				ID:           uuid.New(),
				IsBlacklisted: false,
			},
			expected: true,
		},
		{
			name: "blacklisted user",
			user: &User{
				ID:           uuid.New(),
				IsBlacklisted: true,
			},
			expected: false,
		},
		{
			name: "nil user ID",
			user: &User{
				ID:           uuid.Nil,
				IsBlacklisted: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.IsValid()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
