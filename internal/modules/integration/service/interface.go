package service

import (
	"context"

	"rentos-backend/internal/modules/integration/dto/request"
	"rentos-backend/internal/modules/integration/dto/response"
	"rentos-backend/internal/modules/integration/entity"
)

// APIKeyService manages API key lifecycle.
type APIKeyService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateAPIKey) (*response.APIKeyCreated, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.APIKey, error)
	List(ctx context.Context, tenantID string) ([]response.APIKey, error)
	Revoke(ctx context.Context, id, tenantID string) error
	// ResolveByRawKey validates an incoming raw key and returns the entity
	// (including tenant_id + scopes) for the APIKeyAuthMiddleware.
	ResolveByRawKey(ctx context.Context, rawKey string) (*entity.APIKey, error)
}

// WebhookService manages webhook registrations and dispatch.
type WebhookService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateWebhook) (*response.Webhook, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Webhook, error)
	List(ctx context.Context, tenantID string) ([]response.Webhook, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateWebhook) (*response.Webhook, error)
	Delete(ctx context.Context, id, tenantID string) error
	ListLogs(ctx context.Context, webhookID, tenantID string, page, perPage int) ([]response.WebhookLog, error)
	// Dispatch is called internally by domain event handlers to fan-out
	// webhook delivery for a given event type. It writes webhook_logs and
	// enqueues HTTP delivery (performed by the webhook worker in Phase workers).
	Dispatch(ctx context.Context, req request.DispatchWebhook) error
}
