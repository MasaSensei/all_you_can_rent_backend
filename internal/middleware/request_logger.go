package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// RequestLogger logs one structured line per completed request. Mount it
// after RequestID so the request ID is available to include.
func RequestLogger(base zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		base.Info().
			Str("request_id", GetRequestID(c)).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", c.Response().StatusCode()).
			Dur("latency", time.Since(start)).
			Msg("request handled")

		return err
	}
}
