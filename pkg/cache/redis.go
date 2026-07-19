// Package cache wraps the Redis client construction. It returns the
// concrete *redis.Client rather than introducing an extra interface,
// since go-redis already exposes a stable, mockable surface
// (redis.Cmdable) that later modules can depend on directly.
package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Config holds Redis connection parameters.
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Addr returns the host:port address string.
func (c Config) Addr() string { return fmt.Sprintf("%s:%d", c.Host, c.Port) }

// New creates a Redis client.
func New(cfg Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

// Health pings Redis, bounded by the caller's context deadline.
func Health(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}
