package integration

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/integration/handler"
	"rentos-backend/internal/modules/integration/repository/postgres"
	"rentos-backend/internal/modules/integration/routes"
	"rentos-backend/internal/modules/integration/service"
)

// Module holds the integration module's wired handler and services.
type Module struct {
	handler        *handler.Handler
	apiKeySvc      service.APIKeyService
	webhookSvc     service.WebhookService
}

// New builds the integration module: repositories → services → handler.
func New(c *bootstrap.Container) *Module {
	apiKeyRepo := postgres.NewAPIKeyRepository(
		query("create_api_key.sql"),
		query("find_api_key_by_id.sql"),
		query("find_api_key_by_hash.sql"),
		query("list_api_keys.sql"),
		query("update_api_key_last_used.sql"),
		query("revoke_api_key.sql"),
	)
	webhookRepo := postgres.NewWebhookRepository(
		query("create_webhook.sql"),
		query("find_webhook_by_id.sql"),
		query("list_webhooks.sql"),
		query("list_webhooks_for_event.sql"),
		query("update_webhook.sql"),
		query("delete_webhook.sql"),
	)
	webhookLogRepo := postgres.NewWebhookLogRepository(
		query("create_webhook_log.sql"),
		query("list_webhook_logs.sql"),
	)

	apiKeySvc  := service.NewAPIKeyService(c.DB, apiKeyRepo)
	webhookSvc := service.NewWebhookService(c.DB, webhookRepo, webhookLogRepo)

	h := handler.New(apiKeySvc, webhookSvc, c.Validator)
	return &Module{handler: h, apiKeySvc: apiKeySvc, webhookSvc: webhookSvc}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}

// APIKeyService exposes the service so the APIKeyAuthMiddleware can call
// ResolveByRawKey on every request bearing an API key header.
func (m *Module) APIKeyService() service.APIKeyService {
	return m.apiKeySvc
}

// WebhookService exposes the service so domain event handlers in other
// modules can call Dispatch to fan-out event notifications.
func (m *Module) WebhookService() service.WebhookService {
	return m.webhookSvc
}
