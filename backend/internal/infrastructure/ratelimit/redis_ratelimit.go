package ratelimit

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/response"
)

type RedisRateLimiter struct {
	client  *redis.Client
	rate    int
	burst   int
	logger  *slog.Logger
	enabled bool
}

func NewRedisRateLimiter(client *redis.Client, rate, burst int, logger *slog.Logger) *RedisRateLimiter {
	enabled := true
	if client == nil {
		enabled = false
		logger.Warn("redis client is nil, rate limiter disabled")
	}
	return &RedisRateLimiter{
		client:  client,
		rate:    rate,
		burst:   burst,
		logger:  logger,
		enabled: enabled,
	}
}

func (rl *RedisRateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !rl.enabled {
			next.ServeHTTP(w, r)
			return
		}

		ip := r.RemoteAddr
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			ip = strings.SplitN(fwd, ",", 2)[0]
		}
		if xff := r.Header.Get("X-Real-IP"); xff != "" {
			ip = xff
		}

		key := fmt.Sprintf("ratelimit:%s", ip)
		now := time.Now().Unix()
		window := int64(1)

		pipe := rl.client.Pipeline()
		pipe.ZRemRangeByScore(r.Context(), key, "0", fmt.Sprintf("%d", now-window))
		count := pipe.ZCard(r.Context(), key)
		pipe.ZAdd(r.Context(), key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", now)})
		pipe.Expire(r.Context(), key, 60*time.Second)
		_, err := pipe.Exec(r.Context())

		if err != nil {
			rl.logger.Warn("redis rate limiter error", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		if count.Val() >= int64(rl.burst) {
			rl.logger.Warn("rate limit exceeded", "ip", ip, "count", count.Val())
			w.Header().Set("Retry-After", "1")
			response.JSON(w, http.StatusTooManyRequests, map[string]interface{}{
				"success": false,
				"error":   "rate limit exceeded, try again later",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
