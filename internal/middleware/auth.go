package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"rentos-backend/pkg/jwt"
	apiresponse "rentos-backend/pkg/response"
)

const (
	LocalsUserID    = "user_id"
	LocalsSessionID = "session_id"
	LocalsClaims    = "claims"
)

// Auth returns a middleware that validates the Bearer JWT in the
// Authorization header, then injects user_id, tenant_id, and the full
// *jwt.Claims into fiber.Ctx.Locals for downstream handlers.
//
// Tenant ID from the token always wins over the value set by
// TenantResolver — if a client sends a mismatched X-Tenant-ID header
// we correct it silently, because the token was signed by our auth
// service and is the authoritative source.
func Auth(jwtSvc *jwt.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Get("Authorization")
		if raw == "" {
			return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeUnauthorized, "authorization header is required"))
		}

		parts := strings.SplitN(raw, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeUnauthorized, "authorization header must be: Bearer <token>"))
		}

		claims, err := jwtSvc.ParseAccess(parts[1])
		if err != nil {
			return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeUnauthorized, "invalid or expired token"))
		}

		c.Locals(LocalsUserID, claims.UserID)
		c.Locals(LocalsTenantID, claims.TenantID) // override TenantResolver value
		c.Locals(LocalsClaims, claims)
		return c.Next()
	}
}

// ClaimsFromCtx is a typed helper for handlers and other middleware that
// need the full *jwt.Claims without a type assertion at every call site.
func ClaimsFromCtx(c *fiber.Ctx) *jwt.Claims {
	claims, _ := c.Locals(LocalsClaims).(*jwt.Claims)
	return claims
}
