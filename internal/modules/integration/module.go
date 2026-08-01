package integration

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/integration/handler"
	"rentos-backend/internal/modules/integration/repository/postgres"
	"rentos-backend/internal/modules/integration/routes"
	"rentos-backend/internal/modules/integration/service"
)

type Module struct {
	handler    *handler.Handler
	apiKeySvc  service.APIKeyService
	webhookSvc service.WebhookService
}

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
		query("update_webhook_log_status.sql"),
	)

	apiKeySvc := service.NewAPIKeyService(c.DB, apiKeyRepo)
	webhookSvc := service.NewWebhookService(c.DB, webhookRepo, webhookLogRepo)

	h := handler.New(apiKeySvc, webhookSvc, c.Validator)
	return &Module{handler: h, apiKeySvc: apiKeySvc, webhookSvc: webhookSvc}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}

func (m *Module) APIKeyService() service.APIKeyService   { return m.apiKeySvc }
func (m *Module) WebhookService() service.WebhookService { return m.webhookSvc }
