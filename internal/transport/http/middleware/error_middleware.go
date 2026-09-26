package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/response"
)

// ErrorHandler catches panics and standardizes errors attached to the Gin context
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				response.Error(c, http.StatusInternalServerError, "Internal Server Error")
				c.Abort()
			}
		}()

		c.Next()

		// If any errors were attached to the context (e.g. c.Error(err)),
		// format them into our standard JSON response
		if len(c.Errors) > 0 {
			// In a real production app, you might inspect the error type
			// to return 400 Bad Request vs 500 Internal Server Error.
			// Here we default to 400 for demonstration.
			errMessage := c.Errors.Last().Error()
			response.Error(c, http.StatusBadRequest, errMessage)
			return
		}
	}
}
