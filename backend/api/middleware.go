package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter returns a Gin middleware that uses a token bucket to rate limit requests.
func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}
		c.Next()
	}
}
