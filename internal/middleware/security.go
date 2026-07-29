package middleware

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/pkg/response"
)

// BodySizeLimit rejects requests whose body exceeds maxBytes.
// Fiber has a global BodyLimit in Config, but this lets you set per-group limits.
func BodySizeLimit(maxBytes int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Body()) > maxBytes {
			return response.Error(c, response.NewAppError(
				response.CodeValidation,
				"request body too large",
			))
		}
		return c.Next()
	}
}

// SecurityHeaders adds standard security response headers.
// Mount globally, before route handlers.
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Content-Security-Policy", "default-src 'none'")
		// HSTS only on production — callers should gate on env.
		return c.Next()
	}
}

// RequestID is already in request_id.go; exposed here for reference.

// NoCache sets Cache-Control headers to prevent caching on API responses.
func NoCache() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Set("Pragma", "no-cache")
		return c.Next()
	}
}
