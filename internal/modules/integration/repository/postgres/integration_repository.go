package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos-backend/internal/modules/integration/entity"
	"rentos-backend/internal/modules/integration/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// apiKeyRepository
// ============================================================

type apiKeyRepository struct {
	qCreate         string
	qFindByID       string
	qFindByHash     string
	qList           string
	qUpdateLastUsed string
	qRevoke         string
}

func NewAPIKeyRepository(qCreate, qFindByID, qFindByHash, qList, qUpdateLastUsed, qRevoke string) repository.APIKeyRepository {
	return &apiKeyRepository{
		qCreate: qCreate, qFindByID: qFindByID, qFindByHash: qFindByHash,
		qList: qList, qUpdateLastUsed: qUpdateLastUsed, qRevoke: qRevoke,
	}
}

func (r *apiKeyRepository) Create(ctx context.Context, q database.Querier, k *entity.APIKey) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		k.ID, k.TenantID, k.Name, k.KeyPrefix, k.KeyHash,
		k.Scopes, k.ExpiresAt, k.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("apiKeyRepository.Create: %w", err)
	}
	return nil
}

func (r *apiKeyRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.APIKey, error) {
	var k entity.APIKey
	if err := q.GetContext(ctx, &k, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("apiKeyRepository.FindByID: %w", err)
	}
	return &k, nil
}

func (r *apiKeyRepository) FindByHash(ctx context.Context, q database.Querier, hash string) (*entity.APIKey, error) {
	var k entity.APIKey
	if err := q.GetContext(ctx, &k, r.qFindByHash, hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("apiKeyRepository.FindByHash: %w", err)
	}
	return &k, nil
}

func (r *apiKeyRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.APIKey, error) {
	var out []entity.APIKey
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("apiKeyRepository.List: %w", err)
	}
	return out, nil
}

func (r *apiKeyRepository) UpdateLastUsed(ctx context.Context, q database.Querier, id string) error {
	_, err := q.ExecContext(ctx, r.qUpdateLastUsed, id)
	return err
}

func (r *apiKeyRepository) Revoke(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qRevoke, id, tenantID)
	if err != nil {
		return fmt.Errorf("apiKeyRepository.Revoke: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// webhookRepository
// ============================================================

type webhookRepository struct {
	qCreate       string
	qFindByID     string
	qList         string
	qListForEvent string
	qUpdate       string
	qDelete       string
}

func NewWebhookRepository(qCreate, qFindByID, qList, qListForEvent, qUpdate, qDelete string) repository.WebhookRepository {
	return &webhookRepository{
		qCreate: qCreate, qFindByID: qFindByID, qList: qList,
		qListForEvent: qListForEvent, qUpdate: qUpdate, qDelete: qDelete,
	}
}

func (r *webhookRepository) Create(ctx context.Context, q database.Querier, w *entity.Webhook) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		w.ID, w.TenantID, w.URL, w.Events, w.Secret, w.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("webhookRepository.Create: %w", err)
	}
	return nil
}

func (r *webhookRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Webhook, error) {
	var w entity.Webhook
	if err := q.GetContext(ctx, &w, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("webhookRepository.FindByID: %w", err)
	}
	return &w, nil
}

func (r *webhookRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Webhook, error) {
	var out []entity.Webhook
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("webhookRepository.List: %w", err)
	}
	return out, nil
}

func (r *webhookRepository) ListForEvent(ctx context.Context, q database.Querier, tenantID, eventType string) ([]entity.Webhook, error) {
	var out []entity.Webhook
	if err := q.SelectContext(ctx, &out, r.qListForEvent, tenantID, eventType); err != nil {
		return nil, fmt.Errorf("webhookRepository.ListForEvent: %w", err)
	}
	return out, nil
}

func (r *webhookRepository) Update(ctx context.Context, q database.Querier, w *entity.Webhook) error {
	res, err := q.ExecContext(ctx, r.qUpdate,
		w.ID, w.TenantID, w.URL, w.Events, w.IsActive,
	)
	if err != nil {
		return fmt.Errorf("webhookRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *webhookRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("webhookRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// webhookLogRepository
// ============================================================

type webhookLogRepository struct {
	qCreate        string
	qListByWebhook string
	qUpdateStatus  string
}

func NewWebhookLogRepository(qCreate, qListByWebhook, qUpdateStatus string) repository.WebhookLogRepository {
	return &webhookLogRepository{
		qCreate:        qCreate,
		qListByWebhook: qListByWebhook,
		qUpdateStatus:  qUpdateStatus,
	}
}

func (r *webhookLogRepository) Create(ctx context.Context, q database.Querier, l *entity.WebhookLog) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		l.ID, l.TenantID, l.WebhookID, l.EventType, l.Payload,
		l.ResponseCode, l.ResponseBody,
	)
	if err != nil {
		return fmt.Errorf("webhookLogRepository.Create: %w", err)
	}
	return nil
}

func (r *webhookLogRepository) ListByWebhook(ctx context.Context, q database.Querier, webhookID, tenantID string, limit, offset int) ([]entity.WebhookLog, error) {
	var out []entity.WebhookLog
	if err := q.SelectContext(ctx, &out, r.qListByWebhook, webhookID, tenantID, limit, offset); err != nil {
		return nil, fmt.Errorf("webhookLogRepository.ListByWebhook: %w", err)
	}
	return out, nil
}

func (r *webhookLogRepository) UpdateStatus(ctx context.Context, q database.Querier, id, deliveryStatus string, responseCode *int, responseBody *string) error {
	_, err := q.ExecContext(ctx, r.qUpdateStatus, id, deliveryStatus, responseCode, responseBody)
	if err != nil {
		return fmt.Errorf("webhookLogRepository.UpdateStatus: %w", err)
	}
	return nil
}
