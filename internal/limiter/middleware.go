package limiter

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	Limiter *TokenBucketLimiter
	Max     int
	Rate    float64
}

func NewMiddleware(cacheDB *redis.Client, maxTokens int, refillRate float64) *Middleware {
	return &Middleware{
		Limiter: NewTokenBucketLimiter(cacheDB),
		Max:     maxTokens,
		Rate:    refillRate,
	}
}

func (m *Middleware) TokenBucketHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, tokensLeft, err := m.Limiter.AllowRequest(context.Background(), c.ClientIP(), m.Max, m.Rate)

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

		// Tudo certo, continua para o próximo handler
		c.Next()
	}
}
