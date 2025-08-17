package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery recovers from panics and provides proper error responses
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic with stack trace
		c.Error(fmt.Errorf("panic recovered: %v\n%s", recovered, debug.Stack()))

		// Return a generic error response to the client
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			},
		})
	})
}
