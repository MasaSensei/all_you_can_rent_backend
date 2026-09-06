package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/superadmin/handler"
)

// RegisterPublic mounts routes accessible without auth (register + plan list).
func RegisterPublic(router fiber.Router, h *handler.Handler) {
	// Tenant self-registration
	router.Post("/auth/register", h.RegisterTenant)
	// Public plan listing (for the register page plan selector)
	router.Get("/plans", h.ListPlans)
}

// RegisterSuperAdmin mounts super admin routes — should be behind super admin auth middleware.
func RegisterSuperAdmin(router fiber.Router, h *handler.Handler) {
	router.Post("/login", h.Login)

	// Protected — require super admin JWT
	router.Get("/tenants", h.ListTenants)
	router.Get("/tenants/:id", h.GetTenant)
	router.Patch("/tenants/:id/status", h.UpdateTenantStatus)
	router.Post("/tenants/:id/subscriptions", h.AssignSubscription)
	router.Get("/plans", h.ListPlans)
	router.Get("/stats", h.Stats)
}
