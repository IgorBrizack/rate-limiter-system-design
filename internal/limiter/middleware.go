package limiter

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	TokenBucketLimiter  *TokenBucketLimiter
	LikingBucketLimiter *LeakingBucketLimiter
}

func NewMiddleware(cacheDB *redis.Client) *Middleware {
	return &Middleware{
		TokenBucketLimiter:  NewTokenBucketLimiter(cacheDB),
		LikingBucketLimiter: NewLeakingBucketLimiter(cacheDB),
	}
}

func (m *Middleware) TokenBucketHandler(Max int,
	Rate float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, tokensLeft, err := m.TokenBucketLimiter.AllowRequest(context.Background(), c.ClientIP(), Max, Rate)

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

func (m *Middleware) LeakingBucketHandler(Max int,
	Rate float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, tokensLeft, err := m.LikingBucketLimiter.AllowRequest(context.Background(), c.ClientIP(), Max, Rate)

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
