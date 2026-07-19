package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos/internal/modules/core/dto/request"
	"rentos/internal/modules/core/entity"
	"rentos/internal/modules/core/repository"
	"rentos/pkg/response"
)

const (
	defaultAuditPageSize = 20
	maxAuditPageSize     = 100
)

// ============================================================
// TenantService
// ============================================================

type tenantService struct {
	db   *sqlx.DB
	repo repository.TenantRepository
}

func NewTenantService(db *sqlx.DB, repo repository.TenantRepository) TenantService {
	return &tenantService{db: db, repo: repo}
}

func (s *tenantService) CreateTenant(ctx context.Context, req request.CreateTenant) (*entity.Tenant, error) {
	// Slug uniqueness guard — FindBySlug returning ErrNotFound is the
	// happy path here; any other error is a real failure.
	existing, err := s.repo.FindBySlug(ctx, s.db, req.Slug)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, response.NewAppError(response.CodeConflict, "a tenant with this slug already exists")
	}

	t := &entity.Tenant{
		ID:     uuid.NewString(),
		Name:   req.Name,
		Slug:   req.Slug,
		Domain: req.Domain,
		Plan:   req.Plan,
		Status: "pending",
	}
	if err := s.repo.Create(ctx, s.db, t); err != nil {
		return nil, err
	}
	// Full provisioning (seed default roles/permissions, create owner
	// user) is added once the rbac and auth modules exist in Phase 2.
	return t, nil
}

func (s *tenantService) GetTenantByID(ctx context.Context, id string) (*entity.Tenant, error) {
	t, err := s.repo.FindByID(ctx, s.db, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.NewAppError(response.CodeNotFound, "tenant not found")
		}
		return nil, err
	}
	return t, nil
}

func (s *tenantService) GetTenantBySlug(ctx context.Context, slug string) (*entity.Tenant, error) {
	t, err := s.repo.FindBySlug(ctx, s.db, slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.NewAppError(response.CodeNotFound, "tenant not found")
		}
		return nil, err
	}
	return t, nil
}

// ============================================================
// SettingService
// ============================================================

type settingService struct {
	db          *sqlx.DB
	settings    repository.SettingRepository
	sysSettings repository.SystemSettingRepository
}

func NewSettingService(db *sqlx.DB, settings repository.SettingRepository, sysSettings repository.SystemSettingRepository) SettingService {
	return &settingService{db: db, settings: settings, sysSettings: sysSettings}
}

func (s *settingService) UpsertSetting(ctx context.Context, tenantID string, req request.UpsertSetting, actorID *string) (*entity.Setting, error) {
	setting := &entity.Setting{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Key:       req.Key,
		Value:     req.Value,
		Type:      req.Type,
		UpdatedBy: actorID,
	}
	if err := s.settings.Upsert(ctx, s.db, setting); err != nil {
		return nil, err
	}
	return s.GetSetting(ctx, tenantID, req.Key)
}

func (s *settingService) GetSetting(ctx context.Context, tenantID, key string) (*entity.Setting, error) {
	v, err := s.settings.Find(ctx, s.db, tenantID, key)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.NewAppError(response.CodeNotFound, "setting not found")
		}
		return nil, err
	}
	return v, nil
}

func (s *settingService) ListSettings(ctx context.Context, tenantID string) ([]entity.Setting, error) {
	return s.settings.ListByTenant(ctx, s.db, tenantID)
}

func (s *settingService) UpsertSystemSetting(ctx context.Context, req request.UpsertSystemSetting) (*entity.SystemSetting, error) {
	setting := &entity.SystemSetting{
		ID:    uuid.NewString(),
		Key:   req.Key,
		Value: req.Value,
		Type:  req.Type,
	}
	if err := s.sysSettings.Upsert(ctx, s.db, setting); err != nil {
		return nil, err
	}
	return s.GetSystemSetting(ctx, req.Key)
}

func (s *settingService) GetSystemSetting(ctx context.Context, key string) (*entity.SystemSetting, error) {
	v, err := s.sysSettings.Find(ctx, s.db, key)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.NewAppError(response.CodeNotFound, "system setting not found")
		}
		return nil, err
	}
	return v, nil
}

// ============================================================
// AuditService
// ============================================================

type auditService struct {
	db   *sqlx.DB
	repo repository.AuditLogRepository
}

func NewAuditService(db *sqlx.DB, repo repository.AuditLogRepository) AuditService {
	return &auditService{db: db, repo: repo}
}

func (s *auditService) Write(ctx context.Context, tenantID string, req request.CreateAuditLog) error {
	log := &entity.AuditLog{
		ID:         uuid.NewString(),
		TenantID:   tenantID,
		UserID:     req.UserID,
		EntityType: req.EntityType,
		EntityID:   req.EntityID,
		Action:     req.Action,
		OldValues:  req.OldValues,
		NewValues:  req.NewValues,
		IPAddress:  req.IPAddress,
		UserAgent:  req.UserAgent,
	}
	return s.repo.Create(ctx, s.db, log)
}

func (s *auditService) List(ctx context.Context, tenantID string, filter request.ListAuditLogsFilter) ([]entity.AuditLog, error) {
	perPage := filter.PerPage
	if perPage <= 0 {
		perPage = defaultAuditPageSize
	}
	if perPage > maxAuditPageSize {
		perPage = maxAuditPageSize
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * perPage

	return s.repo.List(ctx, s.db, tenantID, filter.EntityType, filter.EntityID, filter.Action, perPage, offset)
}
