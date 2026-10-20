package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	infraredis "github.com/rmValdez/go-gin-template-v1/internal/infrastructure/redis"
	"github.com/rmValdez/go-gin-template-v1/internal/transport/http/response"
)

// rateLimitScript atomically increments a counter and sets expiry if it's new
var rateLimitScript = redis.NewScript(`
local current = redis.call('INCR', KEYS[1])
if current == 1 then
    redis.call('PEXPIRE', KEYS[1], ARGV[1])
end
return current
`)

// RateLimiter returns a middleware that limits requests by IP address using Redis atomically.
func RateLimiter(prefix string, limit int, window time.Duration) gin.HandlerFunc {
	windowMs := window.Milliseconds()

	return func(c *gin.Context) {
		if infraredis.Client == nil {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s:%s", prefix, clientIP)

		ctx, cancel := context.WithTimeout(c.Request.Context(), 500*time.Millisecond)
		defer cancel()

		result, err := rateLimitScript.Run(ctx, infraredis.Client, []string{key}, windowMs).Result()
		if err != nil {
			// On Redis error or timeout, fail open to avoid service outage
			c.Next()
			return
		}

		count, ok := result.(int64)
		if !ok {
			c.Next()
			return
		}

		remaining := int64(limit) - count
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))

		if count > int64(limit) {
			c.Header("Retry-After", strconv.Itoa(int(window.Seconds())))
			response.Error(c, http.StatusTooManyRequests, "Too many requests. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}
