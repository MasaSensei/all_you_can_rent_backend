package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	integrationSvc "rentos-backend/internal/modules/integration/service"
	"rentos-backend/pkg/response"
)

// APIKeyAuth returns a Fiber middleware that authenticates requests bearing
// an X-API-Key header. On success it injects tenant_id and api_key_id into
// Locals so downstream handlers can read them identically to JWT auth.
//
// Usage in main.go:
//
//	v1.Use(middleware.APIKeyAuth(integration.APIKeyService()))
//
// or mount on a dedicated /v1/external group:
//
//	external := v1.Group("/external", middleware.APIKeyAuth(integration.APIKeyService()))
func APIKeyAuth(svc integrationSvc.APIKeyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := c.Get("X-API-Key")
		if raw == "" {
			// Also accept Bearer token format: "Bearer rnt_..."
			auth := c.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer rnt_") {
				raw = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if raw == "" {
			return response.Error(c, response.NewAppError(response.CodeUnauthorized, "API key required"))
		}

		key, err := svc.ResolveByRawKey(c.Context(), raw)
		if err != nil {
			return response.FromError(c, err)
		}

		// Inject the same locals as JWTAuth so handlers are auth-agnostic.
		c.Locals("tenant_id", key.TenantID)
		c.Locals("api_key_id", key.ID)
		c.Locals("auth_method", "api_key")

		return c.Next()
	}
}
