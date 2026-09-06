package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos-backend/internal/modules/superadmin/dto/response"
	"rentos-backend/internal/modules/superadmin/entity"
	"rentos-backend/internal/modules/superadmin/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// tenantRepository
// ============================================================

type tenantRepository struct {
	qCreate                   string
	qFindByID                 string
	qFindBySlug               string
	qSlugExists               string
	qList                     string
	qUpdateStatus             string
	qUpdateSubscriptionStatus string
	qPlatformStats            string
}

func NewTenantRepository(
	qCreate, qFindByID, qFindBySlug, qSlugExists, qList,
	qUpdateStatus, qUpdateSubscriptionStatus, qPlatformStats string,
) repository.TenantRepository {
	return &tenantRepository{
		qCreate: qCreate, qFindByID: qFindByID, qFindBySlug: qFindBySlug,
		qSlugExists: qSlugExists, qList: qList, qUpdateStatus: qUpdateStatus,
		qUpdateSubscriptionStatus: qUpdateSubscriptionStatus, qPlatformStats: qPlatformStats,
	}
}

func (r *tenantRepository) Create(ctx context.Context, q database.Querier, t *entity.Tenant) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		t.ID, t.Name, t.Slug, t.Email, t.Phone,
		t.Timezone, t.Locale, t.Currency,
		t.SubscriptionStatus, t.TrialEndsAt, t.Status, t.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("tenantRepository.Create: %w", err)
	}
	return nil
}

func (r *tenantRepository) FindByID(ctx context.Context, q database.Querier, id string) (*entity.Tenant, error) {
	var t entity.Tenant
	if err := q.GetContext(ctx, &t, r.qFindByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("tenantRepository.FindByID: %w", err)
	}
	return &t, nil
}

func (r *tenantRepository) FindBySlug(ctx context.Context, q database.Querier, slug string) (*entity.Tenant, error) {
	var t entity.Tenant
	if err := q.GetContext(ctx, &t, r.qFindBySlug, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("tenantRepository.FindBySlug: %w", err)
	}
	return &t, nil
}

func (r *tenantRepository) SlugExists(ctx context.Context, q database.Querier, slug string) (bool, error) {
	var exists bool
	if err := q.GetContext(ctx, &exists, r.qSlugExists, slug); err != nil {
		return false, fmt.Errorf("tenantRepository.SlugExists: %w", err)
	}
	return exists, nil
}

func (r *tenantRepository) List(ctx context.Context, q database.Querier, search, subStatus, status string, limit, offset int) ([]entity.Tenant, error) {
	var out []entity.Tenant
	if err := q.SelectContext(ctx, &out, r.qList, search, subStatus, status, limit, offset); err != nil {
		return nil, fmt.Errorf("tenantRepository.List: %w", err)
	}
	return out, nil
}

func (r *tenantRepository) UpdateStatus(ctx context.Context, q database.Querier, id, status string) error {
	_, err := q.ExecContext(ctx, r.qUpdateStatus, id, status)
	return err
}

func (r *tenantRepository) UpdateSubscriptionStatus(ctx context.Context, q database.Querier, id, status string) error {
	_, err := q.ExecContext(ctx, r.qUpdateSubscriptionStatus, id, status)
	return err
}

func (r *tenantRepository) GetPlatformStats(ctx context.Context, q database.Querier) (*response.PlatformStats, error) {
	var s response.PlatformStats
	if err := q.GetContext(ctx, &s, r.qPlatformStats); err != nil {
		return nil, fmt.Errorf("tenantRepository.GetPlatformStats: %w", err)
	}
	return &s, nil
}

// ============================================================
// subscriptionPlanRepository
// ============================================================

type subscriptionPlanRepository struct {
	qFindBySlug string
	qListActive string
}

func NewSubscriptionPlanRepository(qFindBySlug, qListActive string) repository.SubscriptionPlanRepository {
	return &subscriptionPlanRepository{qFindBySlug: qFindBySlug, qListActive: qListActive}
}

func (r *subscriptionPlanRepository) FindBySlug(ctx context.Context, q database.Querier, slug string) (*entity.SubscriptionPlan, error) {
	var p entity.SubscriptionPlan
	if err := q.GetContext(ctx, &p, r.qFindBySlug, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("subscriptionPlanRepository.FindBySlug: %w", err)
	}
	return &p, nil
}

func (r *subscriptionPlanRepository) ListActive(ctx context.Context, q database.Querier) ([]entity.SubscriptionPlan, error) {
	var out []entity.SubscriptionPlan
	if err := q.SelectContext(ctx, &out, r.qListActive); err != nil {
		return nil, fmt.Errorf("subscriptionPlanRepository.ListActive: %w", err)
	}
	return out, nil
}

// ============================================================
// tenantSubscriptionRepository
// ============================================================

type tenantSubscriptionRepository struct {
	qCreate         string
	qFindActiveByTenant string
}

func NewTenantSubscriptionRepository(qCreate, qFindActiveByTenant string) repository.TenantSubscriptionRepository {
	return &tenantSubscriptionRepository{qCreate: qCreate, qFindActiveByTenant: qFindActiveByTenant}
}

func (r *tenantSubscriptionRepository) Create(ctx context.Context, q database.Querier, s *entity.TenantSubscription) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		s.ID, s.TenantID, s.PlanID, s.BillingCycle,
		s.PricePaid, s.StartsAt, s.EndsAt, s.AutoRenew, s.Status,
	)
	if err != nil {
		return fmt.Errorf("tenantSubscriptionRepository.Create: %w", err)
	}
	return nil
}

func (r *tenantSubscriptionRepository) FindActiveByTenant(ctx context.Context, q database.Querier, tenantID string) (*entity.TenantSubscription, error) {
	var s entity.TenantSubscription
	if err := q.GetContext(ctx, &s, r.qFindActiveByTenant, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("tenantSubscriptionRepository.FindActiveByTenant: %w", err)
	}
	return &s, nil
}

// ============================================================
// superAdminRepository
// ============================================================

type superAdminRepository struct {
	qFindByEmail    string
	qUpdateLastLogin string
}

func NewSuperAdminRepository(qFindByEmail, qUpdateLastLogin string) repository.SuperAdminRepository {
	return &superAdminRepository{qFindByEmail: qFindByEmail, qUpdateLastLogin: qUpdateLastLogin}
}

func (r *superAdminRepository) FindByEmail(ctx context.Context, q database.Querier, email string) (*entity.SuperAdmin, error) {
	var a entity.SuperAdmin
	if err := q.GetContext(ctx, &a, r.qFindByEmail, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("superAdminRepository.FindByEmail: %w", err)
	}
	return &a, nil
}

func (r *superAdminRepository) UpdateLastLogin(ctx context.Context, q database.Querier, id string) error {
	_, err := q.ExecContext(ctx, r.qUpdateLastLogin, id)
	return err
}
