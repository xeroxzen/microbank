package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"microbank/client-service/internal/models"
	"microbank/client-service/internal/services"
)

// UserHandler handles user profile-related HTTP requests
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile retrieves the current user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
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

	// Get user profile
	user, err := h.userService.GetUserByID(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "USER_NOT_FOUND",
				"message": "User not found",
				"details": err.Error(),
			},
		})
		return
	}

	// Return user profile
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"profile": user.ToResponse(),
	})
}

// UpdateProfile updates the current user's profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
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

	// Bind and validate request body
	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid request data",
				"details": err.Error(),
			},
		})
		return
	}

	// Update user profile
	user, err := h.userService.UpdateUserProfile(userUUID, profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "PROFILE_UPDATE_FAILED",
				"message": "Failed to update profile",
				"details": err.Error(),
			},
		})
		return
	}

	// Return updated profile
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": user.ToResponse(),
	})
}
