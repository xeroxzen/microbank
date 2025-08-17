package services

import (
	"fmt"

	"github.com/google/uuid"
	"microbank/client-service/internal/models"
	"microbank/client-service/internal/repository"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// UpdateUserProfile updates a user's profile information
func (s *UserService) UpdateUserProfile(userID uuid.UUID, profile models.UserProfile) (*models.User, error) {
	// Get current user
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update profile fields
	user.Name = profile.Name

	// Save updated user
	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// GetAllUsers retrieves all users (admin only)
func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

// BlacklistUser adds a user to the blacklist (admin only)
func (s *UserService) BlacklistUser(userID uuid.UUID) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Update blacklist status
	if err := s.userRepo.UpdateBlacklistStatus(userID, true); err != nil {
		return fmt.Errorf("failed to blacklist user: %w", err)
	}

	return nil
}

// RemoveFromBlacklist removes a user from the blacklist (admin only)
func (s *UserService) RemoveFromBlacklist(userID uuid.UUID) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Update blacklist status
	if err := s.userRepo.UpdateBlacklistStatus(userID, false); err != nil {
		return fmt.Errorf("failed to remove user from blacklist: %w", err)
	}

	return nil
}

// DeleteUser permanently deletes a user (admin only)
func (s *UserService) DeleteUser(userID uuid.UUID) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Delete user
	if err := s.userRepo.DeleteUser(userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
