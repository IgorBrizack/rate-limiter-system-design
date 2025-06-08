package limiter

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	TokenBucketLimiter  *TokenBucketLimiter
	LikingBucketLimiter *LeakingBucketLimiter
	FixedWindowLimiter  *FixedWindowLimiter
}

func NewMiddleware(cacheDB *redis.Client) *Middleware {
	return &Middleware{
		TokenBucketLimiter:  NewTokenBucketLimiter(cacheDB),
		LikingBucketLimiter: NewLeakingBucketLimiter(cacheDB),
		FixedWindowLimiter:  NewFixedWindowLimiter(cacheDB),
	}
}

func (m *Middleware) TokenBucketHandler(max int,
	rate float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, tokensLeft, err := m.TokenBucketLimiter.AllowRequest(context.Background(), c.ClientIP(), max, rate)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"tokens_left": tokensLeft,
			})
			return
		}

		c.Next()
	}
}

func (m *Middleware) LeakingBucketHandler(max int,
	rate float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, tokensLeft, err := m.LikingBucketLimiter.AllowRequest(context.Background(), c.ClientIP(), max, rate)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"tokens_left": tokensLeft,
			})
			return
		}

		c.Next()
	}
}

func (m *Middleware) FixedWindowHandler(windowSize time.Duration, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		allowed, remaining, err := m.FixedWindowLimiter.AllowRequest(c, key, windowSize, limit)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":         "Rate limit exceeded",
				"requests_left": remaining,
				"retry_after_s": int(windowSize.Seconds()),
			})
			return
		}

		c.Next()
	}
}
