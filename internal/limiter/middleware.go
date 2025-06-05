package limiter

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Middileware struct {
	cacheDB *redis.Client
}

func NewMiddleware(cacheDB *redis.Client) *Middileware {
	return &Middileware{
		cacheDB: cacheDB,
	}
}

func (m *Middileware) TokenBucketAllowRequest(ctx context.Context, redisClient *redis.Client, IPAdress string, maxTokens int, refillRate float64) (bool, int, error) {

	tokenBucketLimiter := NewTokenBucketLimiter(m.cacheDB)

	return tokenBucketLimiter.AllowRequest(ctx, redisClient, IPAdress, maxTokens, refillRate)
}
