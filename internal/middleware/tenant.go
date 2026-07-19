package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	apiresponse "rentos/pkg/response"
)

const (
	LocalsTenantID = "tenant_id"
)

// TenantResolver resolves the current tenant from the request and stores
// its ID in fiber.Ctx.Locals so every downstream handler and service can
// read it via c.Locals(LocalsTenantID).
//
// Resolution order:
//  1. X-Tenant-ID header   (API clients, mobile apps)
//  2. Subdomain             (tenant.rentos.io)
//
// If neither yields a non-empty value the request is rejected with 400.
// The actual lookup of tenant_id → tenant record happens in the auth
// module; here we only extract and forward the identifier so that the
// middleware has zero database dependency and stays fast for every
// request, including unauthenticated ones.
func TenantResolver() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Get("X-Tenant-ID")

		if tenantID == "" {
			tenantID = subdomainFromHost(c.Hostname())
		}

		if tenantID == "" {
			return apiresponse.Error(c,
				apiresponse.NewAppError(apiresponse.CodeValidation, "tenant identifier is required (X-Tenant-ID header or subdomain)"),
			)
		}

		c.Locals(LocalsTenantID, tenantID)
		return c.Next()
	}
}

// subdomainFromHost extracts the leftmost label from a hostname, ignoring
// "www" and single-label hosts (e.g. "localhost").
func subdomainFromHost(host string) string {
	// Strip port if present.
	if idx := strings.LastIndex(host, ":"); idx != -1 {
		host = host[:idx]
	}
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return ""
	}
	sub := parts[0]
	if sub == "www" || sub == "" {
		return ""
	}
	return sub
}
