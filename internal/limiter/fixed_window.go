package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type FixedWindowLimiter struct {
	client *redis.Client
}

func NewFixedWindowLimiter(client *redis.Client) *FixedWindowLimiter {
	return &FixedWindowLimiter{client: client}
}

func (l *FixedWindowLimiter) AllowRequest(ctx context.Context, key string, windowSize time.Duration, limit int) (bool, int, error) {
	now := time.Now()
	windowKey := fmt.Sprintf("%s:%d", key, now.Unix()/int64(windowSize.Seconds()))

	count, err := l.client.Incr(ctx, windowKey).Result()
	if err != nil {
		return false, 0, err
	}

	if count == 1 {
		// Define TTL apenas na primeira requisição da janela
		l.client.Expire(ctx, windowKey, windowSize)
	}

	allowed := int(count) <= limit
	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}

	return allowed, remaining, nil
}
