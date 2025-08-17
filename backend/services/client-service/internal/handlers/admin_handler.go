package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"microbank/client-service/internal/services"
)

// AdminHandler handles administrative HTTP requests
type AdminHandler struct {
	userService *services.UserService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(userService *services.UserService) *AdminHandler {
	return &AdminHandler{
		userService: userService,
	}
}

// GetAllClients retrieves all users (admin only)
func (h *AdminHandler) GetAllClients(c *gin.Context) {
	// Get users
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "FETCH_USERS_FAILED",
				"message": "Failed to fetch users",
				"details": err.Error(),
			},
		})
		return
	}

	// Convert users to response format
	var userResponses []gin.H
	for _, user := range users {
		userResponses = append(userResponses, gin.H{
			"id":             user.ID,
			"email":          user.Email,
			"name":           user.Name,
			"is_blacklisted": user.IsBlacklisted,
			"is_admin":       user.IsAdmin,
			"created_at":     user.CreatedAt,
			"updated_at":     user.UpdatedAt,
		})
	}

	// Return users
	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"users":   userResponses,
		"count":   len(userResponses),
	})
}

// BlacklistClient adds a user to the blacklist (admin only)
func (h *AdminHandler) BlacklistClient(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "Invalid user ID format",
			},
		})
		return
	}

	// Blacklist user
	if err := h.userService.BlacklistUser(userID); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "USER_NOT_FOUND",
					"message": "User not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "BLACKLIST_FAILED",
				"message": "Failed to blacklist user",
				"details": err.Error(),
			},
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User blacklisted successfully",
		"user_id": userID,
	})
}

// RemoveFromBlacklist removes a user from the blacklist (admin only)
func (h *AdminHandler) RemoveFromBlacklist(c *gin.Context) {
	// Get user ID from URL parameter
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_USER_ID",
				"message": "Invalid user ID format",
			},
		})
		return
	}

	// Remove user from blacklist
	if err := h.userService.RemoveFromBlacklist(userID); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "USER_NOT_FOUND",
					"message": "User not found",
				},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "REMOVE_FROM_BLACKLIST_FAILED",
				"message": "Failed to remove user from blacklist",
				"details": err.Error(),
			},
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User removed from blacklist successfully",
		"user_id": userID,
	})
}
