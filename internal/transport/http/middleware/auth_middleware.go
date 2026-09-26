package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/domain/repositories"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/response"
)

// AuthRequired verifies the session token in the HttpOnly cookie
func AuthRequired(sessionRepo repositories.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("session_token")
		if err != nil || token == "" {
			response.Error(c, http.StatusUnauthorized, "Authentication required")
			c.Abort()
			return
		}

		session, err := sessionRepo.FindByToken(token)
		if err != nil || session.ExpiresAt.Before(time.Now()) {
			response.Error(c, http.StatusUnauthorized, "Session expired or invalid")
			c.Abort()
			return
		}

		// Store user ID in context for downstream handlers
		c.Set("userID", session.UserID)
		c.Next()
	}
}
