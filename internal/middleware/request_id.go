package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// HeaderRequestID is the response/request header carrying the correlation ID.
const HeaderRequestID = "X-Request-ID"

const localsRequestID = "request_id"

// RequestID assigns a unique ID to every request, reusing one supplied by
// an upstream proxy if present, and exposes it via fiber.Ctx.Locals and the
// response header for client-side correlation.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get(HeaderRequestID)
		if id == "" {
			id = uuid.NewString()
		}
		c.Locals(localsRequestID, id)
		c.Set(HeaderRequestID, id)
		return c.Next()
	}
}

// GetRequestID retrieves the request ID stored by the RequestID middleware.
func GetRequestID(c *fiber.Ctx) string {
	if id, ok := c.Locals(localsRequestID).(string); ok {
		return id
	}
	return ""
}
