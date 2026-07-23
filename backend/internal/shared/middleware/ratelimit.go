package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/response"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rate     int
	burst    int
	logger   *slog.Logger
}

type visitor struct {
	count    int
	lastSeen time.Time
}

func NewRateLimiter(rate, burst int, logger *slog.Logger) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		logger:   logger,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			ip = strings.SplitN(fwd, ",", 2)[0]
		}
		if xff := r.Header.Get("X-Real-IP"); xff != "" {
			ip = xff
		}

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			rl.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		if time.Since(v.lastSeen) > time.Second {
			v.count = 0
		}
		v.count++
		v.lastSeen = time.Now()
		currentCount := v.count
		rl.mu.Unlock()

		if currentCount > rl.burst {
			if rl.logger != nil {
				rl.logger.Warn("rate limit exceeded",
					"ip", ip,
					"path", r.URL.Path,
					"count", currentCount,
				)
			}
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

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
