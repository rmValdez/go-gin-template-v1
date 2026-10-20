package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodySizeLimit limits the maximum allowed request body size to protect against DoS.
func BodySizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}
