package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"

	apiresponse "rentos/pkg/response"
)

// PermissionResolver is the minimal interface RBACMiddleware needs from
// the rbac module. Using an interface here means middleware never imports
// the rbac package — dependency arrow stays correct.
type PermissionResolver interface {
	GetUserPermissionNames(ctx context.Context, userID, tenantID string) ([]string, error)
}

// Require returns a route-level middleware that checks whether the
// authenticated user holds the given permission (format: "module.action",
// e.g. "booking.create"). Mount it after Auth middleware.
//
// Usage:
//
//	router.Post("/bookings", middleware.Require(rbacModule.UserRoleService(), "booking.create"), handler.Create)
func Require(resolver PermissionResolver, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uID, _ := c.Locals(LocalsUserID).(string)
		tID, _ := c.Locals(LocalsTenantID).(string)

		if uID == "" || tID == "" {
			return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeUnauthorized, "authentication required"))
		}

		names, err := resolver.GetUserPermissionNames(c.Context(), uID, tID)
		if err != nil {
			return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeInternal, "could not resolve permissions"))
		}

		for _, n := range names {
			if n == permission {
				return c.Next()
			}
		}

		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeForbidden, "you do not have permission to perform this action"))
	}
}
