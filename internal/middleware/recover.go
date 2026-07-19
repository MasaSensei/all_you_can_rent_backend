package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
)

// Recover wraps Fiber's built-in recover middleware with structured
// logging of the panic, so a panic anywhere downstream returns a clean
// error response instead of crashing the process.
func Recover(base zerolog.Logger) fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			base.Error().
				Str("request_id", GetRequestID(c)).
				Interface("panic", e).
				Msg("recovered from panic")
		},
	})
}
