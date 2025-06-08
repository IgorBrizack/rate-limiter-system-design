package limiter

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type LeakingBucketLimiter struct {
	cacheDB *redis.Client
}

func NewLeakingBucketLimiter(cacheDB *redis.Client) *LeakingBucketLimiter {
	return &LeakingBucketLimiter{cacheDB: cacheDB}
}

func (l *LeakingBucketLimiter) AllowRequest(ctx context.Context, key string, capacity int, leakRate float64) (bool, int, error) {
	now := float64(time.Now().UnixNano()) / 1e9 // segundos com fração

	pipe := l.cacheDB.TxPipeline()
	lastLeakCmd := pipe.Get(ctx, key+":last")
	sizeCmd := pipe.Get(ctx, key+":size")
	_, err := pipe.Exec(ctx)

	var lastLeak float64
	var currentSize int

	if err != nil && err != redis.Nil {
		return false, 0, err
	}

	// Parse valores
	if lastLeakStr, err := lastLeakCmd.Result(); err == nil {
		lastLeak, _ = strconv.ParseFloat(lastLeakStr, 64)
	}
	if sizeStr, err := sizeCmd.Result(); err == nil {
		currentSize, _ = strconv.Atoi(sizeStr)
	}

	// Calcular quanto "vazou"
	elapsed := now - lastLeak
	leaked := int(elapsed * leakRate)

	if leaked > 0 {
		currentSize -= leaked
		if currentSize < 0 {
			currentSize = 0
		}
		lastLeak = now
	}

	allowed := currentSize < capacity
	if allowed {
		currentSize++
	}

	// Atualizar Redis
	pipe = l.cacheDB.TxPipeline()
	pipe.Set(ctx, key+":size", currentSize, 10*time.Minute)
	pipe.Set(ctx, key+":last", lastLeak, 10*time.Minute)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, currentSize, err
	}

	return allowed, capacity - currentSize, nil
}
