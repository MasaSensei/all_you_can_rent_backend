// Package routes registers the core module's HTTP routes onto a Fiber
// router group. Kept separate from handler so route paths/grouping can
// change without touching handler logic.
package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos/internal/modules/core/handler"
)

// Register mounts the core module's routes under the given router
// (expected to already be scoped to /api/v1).
func Register(router fiber.Router, h *handler.Handler) {
	tenants := router.Group("/tenants")
	tenants.Post("/", h.CreateTenant)
	tenants.Get("/:id", h.GetTenant)

	settings := router.Group("/settings")
	settings.Get("/", h.ListSettings)
	settings.Get("/:key", h.GetSetting)
	settings.Put("/:key", h.UpsertSetting)

	systemSettings := router.Group("/system-settings")
	systemSettings.Get("/:key", h.GetSystemSetting)
	systemSettings.Put("/:key", h.UpsertSystemSetting)

	router.Get("/audit-logs", h.ListAuditLogs)
}
