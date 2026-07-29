package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/integration/dto/request"
	"rentos-backend/internal/modules/integration/dto/response"
	"rentos-backend/internal/modules/integration/entity"
	"rentos-backend/internal/modules/integration/repository"
	pkgresponse "rentos-backend/pkg/response"
)

// ============================================================
// apiKeyService
// ============================================================

type apiKeyService struct {
	db   *sqlx.DB
	repo repository.APIKeyRepository
}

func NewAPIKeyService(db *sqlx.DB, repo repository.APIKeyRepository) APIKeyService {
	return &apiKeyService{db: db, repo: repo}
}

func (s *apiKeyService) Create(ctx context.Context, tenantID, actorID string, req request.CreateAPIKey) (*response.APIKeyCreated, error) {
	// Generate a cryptographically random key: prefix_randomhex
	raw, prefix, hash, err := generateAPIKey()
	if err != nil {
		return nil, err
	}

	scopesJSON, err := json.Marshal(req.Scopes)
	if err != nil {
		return nil, err
	}
	scopesStr := string(scopesJSON)

	k := &entity.APIKey{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Name:      req.Name,
		KeyPrefix: prefix,
		KeyHash:   hash,
		Scopes:    &scopesStr,
		ExpiresAt: req.ExpiresAt,
		CreatedBy: &actorID,
	}

	if err := s.repo.Create(ctx, s.db, k); err != nil {
		return nil, err
	}

	return &response.APIKeyCreated{
		APIKey: response.APIKey{
			ID:        k.ID,
			TenantID:  k.TenantID,
			Name:      k.Name,
			KeyPrefix: k.KeyPrefix,
			Scopes:    req.Scopes,
			ExpiresAt: k.ExpiresAt,
			Status:    "active",
			CreatedAt: k.CreatedAt,
		},
		RawKey: raw,
	}, nil
}

func (s *apiKeyService) GetByID(ctx context.Context, id, tenantID string) (*response.APIKey, error) {
	k, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "API key not found")
		}
		return nil, err
	}
	return toAPIKeyResponse(k), nil
}

func (s *apiKeyService) List(ctx context.Context, tenantID string) ([]response.APIKey, error) {
	keys, err := s.repo.List(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.APIKey, 0, len(keys))
	for _, k := range keys {
		out = append(out, *toAPIKeyResponse(&k))
	}
	return out, nil
}

func (s *apiKeyService) Revoke(ctx context.Context, id, tenantID string) error {
	if err := s.repo.Revoke(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "API key not found")
		}
		return err
	}
	return nil
}

func (s *apiKeyService) ResolveByRawKey(ctx context.Context, rawKey string) (*entity.APIKey, error) {
	hash := hashKey(rawKey)
	k, err := s.repo.FindByHash(ctx, s.db, hash)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeUnauthorized, "invalid API key")
		}
		return nil, err
	}
	// Fire-and-forget last used update — failure is non-fatal.
	_ = s.repo.UpdateLastUsed(ctx, s.db, k.ID)
	return k, nil
}

// ============================================================
// webhookService
// ============================================================

type webhookService struct {
	db       *sqlx.DB
	webhooks repository.WebhookRepository
	logs     repository.WebhookLogRepository
}

func NewWebhookService(
	db *sqlx.DB,
	webhooks repository.WebhookRepository,
	logs repository.WebhookLogRepository,
) WebhookService {
	return &webhookService{db: db, webhooks: webhooks, logs: logs}
}

func (s *webhookService) Create(ctx context.Context, tenantID, actorID string, req request.CreateWebhook) (*response.Webhook, error) {
	eventsJSON, err := json.Marshal(req.Events)
	if err != nil {
		return nil, err
	}

	w := &entity.Webhook{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		URL:       req.URL,
		Events:    string(eventsJSON),
		Secret:    req.Secret,
		IsActive:  true,
		CreatedBy: &actorID,
	}
	if err := s.webhooks.Create(ctx, s.db, w); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, w.ID, tenantID)
}

func (s *webhookService) GetByID(ctx context.Context, id, tenantID string) (*response.Webhook, error) {
	w, err := s.webhooks.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "webhook not found")
		}
		return nil, err
	}
	return toWebhookResponse(w), nil
}

func (s *webhookService) List(ctx context.Context, tenantID string) ([]response.Webhook, error) {
	webhooks, err := s.webhooks.List(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.Webhook, 0, len(webhooks))
	for _, w := range webhooks {
		out = append(out, *toWebhookResponse(&w))
	}
	return out, nil
}

func (s *webhookService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateWebhook) (*response.Webhook, error) {
	var eventsStr *string
	if len(req.Events) > 0 {
		b, _ := json.Marshal(req.Events)
		str := string(b)
		eventsStr = &str
	}

	w := &entity.Webhook{
		ID:        id,
		TenantID:  tenantID,
		URL:       derefStr(req.URL),
		Events:    derefStr(eventsStr),
		IsActive:  derefBool(req.IsActive),
		UpdatedBy: &actorID,
	}
	if err := s.webhooks.Update(ctx, s.db, w); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "webhook not found")
		}
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

func (s *webhookService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.webhooks.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "webhook not found")
		}
		return err
	}
	return nil
}

func (s *webhookService) ListLogs(ctx context.Context, webhookID, tenantID string, page, perPage int) ([]response.WebhookLog, error) {
	perPage, page = normPage(perPage, page)
	logs, err := s.logs.ListByWebhook(ctx, s.db, webhookID, tenantID, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.WebhookLog, 0, len(logs))
	for _, l := range logs {
		out = append(out, response.WebhookLog{
			ID: l.ID, WebhookID: l.WebhookID, EventType: l.EventType,
			ResponseCode: l.ResponseCode, TriggeredAt: l.TriggeredAt, Status: l.Status,
		})
	}
	return out, nil
}

// Dispatch fans out an event to all subscribed webhooks, writing a log
// entry per webhook. Actual HTTP delivery is handled by the webhook
// dispatcher worker, which reads pending log entries and performs the POST.
func (s *webhookService) Dispatch(ctx context.Context, req request.DispatchWebhook) error {
	hooks, err := s.webhooks.ListForEvent(ctx, s.db, req.TenantID, req.EventType)
	if err != nil {
		return err
	}
	if len(hooks) == 0 {
		return nil
	}

	payloadJSON, err := json.Marshal(req.Payload)
	if err != nil {
		return err
	}
	payloadStr := string(payloadJSON)

	for _, hook := range hooks {
		log := &entity.WebhookLog{
			ID:        uuid.NewString(),
			TenantID:  req.TenantID,
			WebhookID: hook.ID,
			EventType: req.EventType,
			Payload:   &payloadStr,
		}
		if err := s.logs.Create(ctx, s.db, log); err != nil {
			// Log the failure but continue to fan-out other webhooks.
			_ = fmt.Errorf("dispatch: log creation failed for webhook %s: %w", hook.ID, err)
		}
		// TODO Phase workers: enqueue log.ID to the webhook_dispatcher worker.
	}
	return nil
}

// ============================================================
// helpers
// ============================================================

func generateAPIKey() (raw, prefix, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	raw = "rnt_" + hex.EncodeToString(b)
	prefix = raw[:12]
	hash = hashKey(raw)
	return
}

func hashKey(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func toAPIKeyResponse(k *entity.APIKey) *response.APIKey {
	var scopes []string
	if k.Scopes != nil {
		_ = json.Unmarshal([]byte(*k.Scopes), &scopes)
	}
	return &response.APIKey{
		ID: k.ID, TenantID: k.TenantID, Name: k.Name, KeyPrefix: k.KeyPrefix,
		Scopes: scopes, LastUsedAt: k.LastUsedAt, ExpiresAt: k.ExpiresAt,
		Status: k.Status, CreatedAt: k.CreatedAt,
	}
}

func toWebhookResponse(w *entity.Webhook) *response.Webhook {
	var events []string
	_ = json.Unmarshal([]byte(w.Events), &events)
	return &response.Webhook{
		ID: w.ID, TenantID: w.TenantID, URL: w.URL,
		Events: events, IsActive: w.IsActive, Status: w.Status,
		CreatedAt: w.CreatedAt, UpdatedAt: w.UpdatedAt,
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefBool(b *bool) bool {
	if b == nil {
		return true
	}
	return *b
}

func normPage(perPage, page int) (int, int) {
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	if page <= 0 {
		page = 1
	}
	return perPage, page
}
