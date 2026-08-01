package repository

import (
	"context"
	"errors"

	"rentos-backend/internal/modules/integration/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

type APIKeyRepository interface {
	Create(ctx context.Context, q database.Querier, k *entity.APIKey) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.APIKey, error)
	FindByHash(ctx context.Context, q database.Querier, hash string) (*entity.APIKey, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.APIKey, error)
	UpdateLastUsed(ctx context.Context, q database.Querier, id string) error
	Revoke(ctx context.Context, q database.Querier, id, tenantID string) error
}

type WebhookRepository interface {
	Create(ctx context.Context, q database.Querier, w *entity.Webhook) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Webhook, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Webhook, error)
	ListForEvent(ctx context.Context, q database.Querier, tenantID, eventType string) ([]entity.Webhook, error)
	Update(ctx context.Context, q database.Querier, w *entity.Webhook) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

type WebhookLogRepository interface {
	Create(ctx context.Context, q database.Querier, l *entity.WebhookLog) error
	ListByWebhook(ctx context.Context, q database.Querier, webhookID, tenantID string, limit, offset int) ([]entity.WebhookLog, error)
	// UpdateStatus dipanggil oleh webhook dispatcher worker setelah HTTP delivery.
	UpdateStatus(ctx context.Context, q database.Querier, id, deliveryStatus string, responseCode *int, responseBody *string) error
}
