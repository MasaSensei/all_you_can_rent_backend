package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/rbac/handler"
)

// Register mounts RBAC routes under /api/v1.
func Register(router fiber.Router, h *handler.Handler) {
	roles := router.Group("/roles")
	roles.Post("/", h.CreateRole)
	roles.Get("/", h.ListRoles)
	roles.Get("/:id", h.GetRole)
	roles.Put("/:id", h.UpdateRole)
	roles.Delete("/:id", h.DeleteRole)
	roles.Put("/:id/permissions", h.SyncPermissions)

	router.Get("/permissions", h.ListPermissions)

	router.Post("/users/:user_id/roles", h.AssignRole)
	router.Delete("/users/:user_id/roles/:role_id", h.RevokeRole)
	router.Get("/users/:user_id/roles", h.ListUserRoles)
}
