package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure (for backward compatibility)
type Claims struct {
	UserID        string `json:"user_id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	IsAdmin       bool   `json:"is_admin"`
	IsBlacklisted bool   `json:"is_blacklisted"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens and extracts user information
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "MISSING_TOKEN",
					"message": "Authorization header is required",
				},
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN_FORMAT",
					"message": "Token must be in format: Bearer <token>",
				},
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate the token
		claims, err := parseAndValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
					"details": err.Error(),
				},
			})
			c.Abort()
			return
		}

		// Check if user is blacklisted
		if claims.IsBlacklisted {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "USER_BLACKLISTED",
					"message": "User account has been suspended",
				},
			})
			c.Abort()
			return
		}

		// Store user information in context for handlers to use
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("is_blacklisted", claims.IsBlacklisted)

		c.Next()
	}
}

// parseAndValidateToken parses and validates a JWT token using MapClaims
func parseAndValidateToken(tokenString string) (*Claims, error) {
	// Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable not set")
	}

	// Parse token using MapClaims
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

	// Extract claims from MapClaims
	if mapClaims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Convert MapClaims to our Claims struct
		claims := &Claims{}
		
		// Extract user_id (required)
		if userID, exists := mapClaims["user_id"]; exists {
			if userIDStr, ok := userID.(string); ok {
				claims.UserID = userIDStr
			} else {
				return nil, fmt.Errorf("invalid user_id type in token")
			}
		} else {
			return nil, fmt.Errorf("user_id not found in token")
		}

		// Extract email (optional)
		if email, exists := mapClaims["email"]; exists {
			if emailStr, ok := email.(string); ok {
				claims.Email = emailStr
			}
		}

		// Extract name (optional)
		if name, exists := mapClaims["name"]; exists {
			if nameStr, ok := name.(string); ok {
				claims.Name = nameStr
			}
		}

		// Extract is_admin (optional, default to false)
		if isAdmin, exists := mapClaims["is_admin"]; exists {
			if isAdminBool, ok := isAdmin.(bool); ok {
				claims.IsAdmin = isAdminBool
			}
		}

		// Extract is_blacklisted (optional, default to false)
		if isBlacklisted, exists := mapClaims["is_blacklisted"]; exists {
			if isBlacklistedBool, ok := isBlacklisted.(bool); ok {
				claims.IsBlacklisted = isBlacklistedBool
			}
		}

		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
