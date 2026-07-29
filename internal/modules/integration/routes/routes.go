package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/integration/handler"
)

// Register mounts all integration routes under /api/v1.
func Register(router fiber.Router, h *handler.Handler) {
	apiKeys := router.Group("/api-keys")
	apiKeys.Post("/", h.CreateAPIKey)
	apiKeys.Get("/", h.ListAPIKeys)
	apiKeys.Get("/:id", h.GetAPIKey)
	apiKeys.Delete("/:id", h.RevokeAPIKey)

	webhooks := router.Group("/webhooks")
	webhooks.Post("/", h.CreateWebhook)
	webhooks.Get("/", h.ListWebhooks)
	webhooks.Get("/:id", h.GetWebhook)
	webhooks.Put("/:id", h.UpdateWebhook)
	webhooks.Delete("/:id", h.DeleteWebhook)
	webhooks.Get("/:id/logs", h.ListWebhookLogs)
}
