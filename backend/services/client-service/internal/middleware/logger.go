package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger provides structured logging for HTTP requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Generate request ID for correlation
		requestID := uuid.New().String()
		
		// Store request ID in context for handlers to use
		if param.Keys == nil {
			param.Keys = make(map[string]interface{})
		}
		param.Keys["request_id"] = requestID

		// Format the log entry
		return fmt.Sprintf("[%s] %s | %s | %d | %s | %s | %s | %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Request.UserAgent(),
			requestID,
		)
	})
}
