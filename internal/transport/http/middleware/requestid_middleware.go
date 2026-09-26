package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a unique UUID to every request for tracing and logging
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Set it in the Gin context so handlers can access it
		c.Set("requestID", requestID)

		// Also attach it to the response header
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
