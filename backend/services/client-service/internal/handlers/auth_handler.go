package handlers

import (
	"net/http"

	"microbank/client-service/internal/models"
	"microbank/client-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var registration models.UserRegistration

	// Bind and validate request body
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
				"details": err.Error(),
			},
		})
		return
	}

	// Register user
	user, err := h.authService.RegisterUser(registration)
	if err != nil {
		// Check for specific error types
		if err.Error() == "user with email "+registration.Email+" already exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": gin.H{
					"code":    "USER_EXISTS",
					"message": "User with this email already exists",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "REGISTRATION_FAILED",
				"message": "Failed to register user",
				"details": err.Error(),
			},
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user.ToResponse(),
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var login models.UserLogin

	// Bind and validate request body
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
				"details": err.Error(),
			},
		})
		return
	}

	// Authenticate user
	user, accessToken, refreshToken, err := h.authService.LoginUser(login)
	if err != nil {
		// Check for specific error types
		if err.Error() == "invalid credentials" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_CREDENTIALS",
					"message": "Invalid email or password",
				},
			})
			return
		}

		if err.Error() == "account has been suspended" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "ACCOUNT_SUSPENDED",
					"message": "Your account has been suspended",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "LOGIN_FAILED",
				"message": "Failed to authenticate user",
				"details": err.Error(),
			},
		})
		return
	}

	// Return success response with tokens
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user.ToResponse(),
		"tokens": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
		},
	})
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	// Bind and validate request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
				"details": err.Error(),
			},
		})
		return
	}

	// Refresh token
	accessToken, err := h.authService.RefreshToken(request.RefreshToken)
	if err != nil {
		// Check for specific error types
		if err.Error() == "invalid refresh token" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_REFRESH_TOKEN",
					"message": "Invalid refresh token",
				},
			})
			return
		}

		if err.Error() == "refresh token expired" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "REFRESH_TOKEN_EXPIRED",
					"message": "Refresh token has expired",
				},
			})
			return
		}

		if err.Error() == "account has been suspended" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "ACCOUNT_SUSPENDED",
					"message": "Your account has been suspended",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "TOKEN_REFRESH_FAILED",
				"message": "Failed to refresh token",
				"details": err.Error(),
			},
		})
		return
	}

	// Return new access token
	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"tokens": gin.H{
			"access_token": accessToken,
			"token_type":   "Bearer",
		},
	})
}

// ValidateToken validates the current access token
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	// Get user information from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "User information not found in context",
			},
		})
		return
	}

	// Parse user ID
	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "Invalid user ID format",
			},
		})
		return
	}


        // Return user information from context
        email, _ := c.Get("email")
        name, _ := c.Get("name")
        isAdmin, _ := c.Get("is_admin")
        c.JSON(http.StatusOK, gin.H{
                "message": "Token is valid",
                "user": gin.H{
                        "id": userUUID,
                        "email": email,
                        "name": name,
                        "is_admin": isAdmin,
                },
        })
}
