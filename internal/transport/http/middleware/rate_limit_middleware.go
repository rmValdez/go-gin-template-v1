package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/response"
	infraredis "github.com/rmValdez/go-gin-template-v1/internal/infrastructure/redis"
)

// RateLimiter returns a middleware that limits requests by IP address using Redis.
// limit is the maximum number of requests allowed within the window.
func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If Redis is not connected (e.g., during tests or if it failed), skip rate limiting
		if infraredis.Client == nil {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		key := "rate_limit:" + clientIP

		ctx := context.Background()

		// Increment the counter for this IP
		count, err := infraredis.Client.Incr(ctx, key).Result()
		if err != nil {
			// On Redis error, just log and allow to prevent blocking users due to cache failure
			c.Next()
			return
		}

		// Set expiration on the very first request in the window
		if count == 1 {
			infraredis.Client.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			response.Error(c, http.StatusTooManyRequests, "Too many requests. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}
