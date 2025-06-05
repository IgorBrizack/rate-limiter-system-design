package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const tokenBucketLua = `
local key = KEYS[1]
local maxTokens = tonumber(ARGV[1])
local refillRate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

local bucket = redis.call("HMGET", key, "tokens", "timestamp")
local tokens = tonumber(bucket[1]) or maxTokens
local lastRefill = tonumber(bucket[2]) or now

local delta = math.max(0, now - lastRefill)
local refill = math.floor(delta * refillRate)
tokens = math.min(maxTokens, tokens + refill)

if tokens < requested then
    return {0, tokens}
else
    tokens = tokens - requested
    redis.call("HMSET", key, "tokens", tokens, "timestamp", now)
    redis.call("EXPIRE", key, 60)
    return {1, tokens}
end
`

type TokenBucketLimiter struct {
	cacheDB *redis.Client
}

func NewTokenBucketLimiter(cacheDB *redis.Client) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		cacheDB: cacheDB,
	}
}

func (l *TokenBucketLimiter) AllowRequest(ctx context.Context, redisClient *redis.Client, IPAdress string, maxTokens int, refillRate float64) (bool, int, error) {
	key := fmt.Sprintf("token_bucket:%s", IPAdress)
	now := time.Now().Unix()

	const consume = 1

	res, err := redisClient.Eval(ctx, tokenBucketLua,
		[]string{key},
		maxTokens,
		refillRate,
		now,
		consume,
	).Result()

	if err != nil {
		return false, 0, err
	}

	values := res.([]interface{})
	allowed := values[0].(int64) == 1
	tokensLeft := int(values[1].(int64))

	return allowed, tokensLeft, nil
}
