package services

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"microbank/client-service/internal/models"
	"microbank/client-service/internal/repository"
)

// AuthService handles authentication-related business logic
type AuthService struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

// RegisterUser handles user registration
func (s *AuthService) RegisterUser(registration models.UserRegistration) (*models.User, error) {
	// Check if user already exists
	exists, err := s.userRepo.UserExists(registration.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("user with email %s already exists", registration.Email)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		ID:           uuid.New(),
		Email:        registration.Email,
		Name:         registration.Name,
		PasswordHash: string(hashedPassword),
		IsBlacklisted: false,
		IsAdmin:      false,
	}

	// Save user to database
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// LoginUser handles user authentication
func (s *AuthService) LoginUser(login models.UserLogin) (*models.User, string, string, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(login.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	// Check if user is blacklisted
	if user.IsBlacklisted {
		return nil, "", "", fmt.Errorf("account has been suspended")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(login.Password)); err != nil {
		return nil, "", "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

// RefreshToken generates a new access token using a refresh token
func (s *AuthService) RefreshToken(refreshTokenString string) (string, error) {
	// Validate refresh token
	refreshToken, err := s.refreshTokenRepo.GetByToken(refreshTokenString)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if refresh token is expired
	if time.Now().After(refreshToken.ExpiresAt) {
		return "", fmt.Errorf("refresh token expired")
	}

	// Get user
	user, err := s.userRepo.GetUserByID(refreshToken.UserID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	// Check if user is blacklisted
	if user.IsBlacklisted {
		return "", fmt.Errorf("account has been suspended")
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return accessToken, nil
}

// ValidateToken validates an access token and returns user information
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	// Parse and validate the token
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract user ID from claims map
	userIDStr, ok := (*claims)["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	// Get user from database to ensure data is current
	user, err := s.userRepo.GetUserByID(uuid.MustParse(userIDStr))
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is blacklisted
	if user.IsBlacklisted {
		return nil, fmt.Errorf("account has been suspended")
	}

	return user, nil
}

// generateAccessToken creates a new JWT access token
func (s *AuthService) generateAccessToken(user *models.User) (string, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}

	// Create claims
	claims := jwt.MapClaims{
		"user_id":        user.ID.String(),
		"email":          user.Email,
		"name":           user.Name,
		"is_admin":       user.IsAdmin,
		"is_blacklisted": user.IsBlacklisted,
		"exp":            time.Now().Add(15 * time.Minute).Unix(), // 15 minutes expiry
		"iat":            time.Now().Unix(),
		"type":           "access",
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}



	return tokenString, nil
}

// generateRefreshToken creates a new refresh token
func (s *AuthService) generateRefreshToken(userID uuid.UUID) (string, error) {
	// Generate a random refresh token
	refreshToken := uuid.New().String()

	// Create refresh token record
	refreshTokenRecord := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: refreshToken, // In production, hash this token
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days expiry
	}

	// Save refresh token to database
	if err := s.refreshTokenRepo.Create(refreshTokenRecord); err != nil {
		return "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return refreshToken, nil
}

// parseToken parses and validates a JWT token
func (s *AuthService) parseToken(tokenString string) (*jwt.MapClaims, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable not set")
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
