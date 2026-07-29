package repository

import (
	"context"
	"errors"

	"rentos-backend/internal/modules/integration/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

// APIKeyRepository manages the api_keys table.
type APIKeyRepository interface {
	Create(ctx context.Context, q database.Querier, k *entity.APIKey) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.APIKey, error)
	// FindByHash resolves an incoming API key by its SHA-256 hash.
	// Used by APIKeyAuthMiddleware on every authenticated request.
	FindByHash(ctx context.Context, q database.Querier, hash string) (*entity.APIKey, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.APIKey, error)
	UpdateLastUsed(ctx context.Context, q database.Querier, id string) error
	Revoke(ctx context.Context, q database.Querier, id, tenantID string) error
}

// WebhookRepository manages the webhooks table.
type WebhookRepository interface {
	Create(ctx context.Context, q database.Querier, w *entity.Webhook) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Webhook, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Webhook, error)
	// ListForEvent returns active webhooks subscribed to a given event type.
	// Used by the webhook dispatcher to fan-out delivery.
	ListForEvent(ctx context.Context, q database.Querier, tenantID, eventType string) ([]entity.Webhook, error)
	Update(ctx context.Context, q database.Querier, w *entity.Webhook) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// WebhookLogRepository manages the webhook_logs table.
type WebhookLogRepository interface {
	Create(ctx context.Context, q database.Querier, l *entity.WebhookLog) error
	ListByWebhook(ctx context.Context, q database.Querier, webhookID, tenantID string, limit, offset int) ([]entity.WebhookLog, error)
	UpdateStatus(ctx context.Context, q database.Querier, logID, status string, errorMsg *string) error
}
