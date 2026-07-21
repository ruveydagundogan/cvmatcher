package cache

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/config"
)

func NewRedisClient(ctx context.Context, cfg config.RedisConfig, logger *slog.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     "",
		DB:           0,
		PoolSize:     50,
		MinIdleConns: 10,
		PoolTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Warn("redis unavailable, running without cache", "error", err)
		return client, nil
	}

	logger.Info("redis connected", "host", cfg.Host, "port", cfg.Port)
	return client, nil
}
