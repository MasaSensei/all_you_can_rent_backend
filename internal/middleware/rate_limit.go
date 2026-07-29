package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"

	"rentos-backend/pkg/response"
)

// ---- Rate Limiter ----

// bucket tracks request count within the current window for one client.
type bucket struct {
	count    int
	resetAt  time.Time
	mu       sync.Mutex
}

// rateLimiter holds per-client buckets in memory.
type rateLimiter struct {
	buckets  sync.Map
	limit    int
	window   time.Duration
}

// RateLimitConfig configures the rate limiter middleware.
type RateLimitConfig struct {
	// Limit is max requests allowed per Window. Default: 100.
	Limit int
	// Window is the rolling window duration. Default: 1 minute.
	Window time.Duration
	// KeyFunc extracts the rate-limit key from the request.
	// Default: client IP.
	KeyFunc func(*fiber.Ctx) string
}

// RateLimit returns a Fiber middleware that enforces per-client request limits.
// Uses an in-memory sliding window. For multi-instance deployments, swap the
// in-memory map for a Redis-backed counter.
func RateLimit(cfg RateLimitConfig) fiber.Handler {
	if cfg.Limit <= 0 {
		cfg.Limit = 100
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if cfg.KeyFunc == nil {
		cfg.KeyFunc = func(c *fiber.Ctx) string { return c.IP() }
	}

	rl := &rateLimiter{limit: cfg.Limit, window: cfg.Window}

	// Background cleanup: sweep expired buckets every window period.
	go func() {
		ticker := time.NewTicker(cfg.Window)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			rl.buckets.Range(func(k, v any) bool {
				b := v.(*bucket)
				b.mu.Lock()
				if now.After(b.resetAt) {
					rl.buckets.Delete(k)
				}
				b.mu.Unlock()
				return true
			})
		}
	}()

	return func(c *fiber.Ctx) error {
		key := cfg.KeyFunc(c)
		now := time.Now()

		val, _ := rl.buckets.LoadOrStore(key, &bucket{resetAt: now.Add(cfg.Window)})
		b := val.(*bucket)

		b.mu.Lock()
		defer b.mu.Unlock()

		// Reset window if expired.
		if now.After(b.resetAt) {
			b.count = 0
			b.resetAt = now.Add(cfg.Window)
		}

		b.count++
		if b.count > rl.limit {
			c.Set("Retry-After", b.resetAt.UTC().Format(time.RFC1123))
			return response.Error(c, response.NewAppError(response.CodeTooManyRequests, "rate limit exceeded"))
		}
		return c.Next()
	}
}

// ---- Presets ----

// RateLimitStrict is for auth endpoints (login, register) — 10 req/min.
func RateLimitStrict() fiber.Handler {
	return RateLimit(RateLimitConfig{Limit: 10, Window: time.Minute})
}

// RateLimitDefault is for general API endpoints — 200 req/min.
func RateLimitDefault() fiber.Handler {
	return RateLimit(RateLimitConfig{Limit: 200, Window: time.Minute})
}

// RateLimitEvents is for analytics event ingestion — 1000 req/min.
func RateLimitEvents() fiber.Handler {
	return RateLimit(RateLimitConfig{Limit: 1000, Window: time.Minute})
}
