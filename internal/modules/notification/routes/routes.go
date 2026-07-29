package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/notification/handler"
)

// Register mounts all notification routes under /api/v1.
func Register(router fiber.Router, h *handler.Handler) {
	// Templates (admin-managed)
	templates := router.Group("/notification-templates")
	templates.Post("/", h.CreateTemplate)
	templates.Get("/", h.ListTemplates)
	templates.Get("/:id", h.GetTemplate)
	templates.Put("/:id", h.UpdateTemplate)
	templates.Delete("/:id", h.DeleteTemplate)

	// In-app notification feed (user-facing)
	notifs := router.Group("/notifications")
	notifs.Post("/", h.CreateInApp)
	notifs.Get("/", h.ListNotifications)
	notifs.Patch("/:id/read", h.MarkRead)
	notifs.Patch("/read-all", h.MarkAllRead)

	// Delivery audit logs (admin-facing)
	router.Get("/notification-logs", h.ListLogs)
}
