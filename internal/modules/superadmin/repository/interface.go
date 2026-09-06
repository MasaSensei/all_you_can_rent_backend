package repository

import (
	"context"
	"errors"

	"rentos-backend/internal/modules/superadmin/dto/response"
	"rentos-backend/internal/modules/superadmin/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

type TenantRepository interface {
	Create(ctx context.Context, q database.Querier, t *entity.Tenant) error
	FindByID(ctx context.Context, q database.Querier, id string) (*entity.Tenant, error)
	FindBySlug(ctx context.Context, q database.Querier, slug string) (*entity.Tenant, error)
	SlugExists(ctx context.Context, q database.Querier, slug string) (bool, error)
	List(ctx context.Context, q database.Querier, search, subscriptionStatus, status string, limit, offset int) ([]entity.Tenant, error)
	UpdateStatus(ctx context.Context, q database.Querier, id, status string) error
	UpdateSubscriptionStatus(ctx context.Context, q database.Querier, id, status string) error
	GetPlatformStats(ctx context.Context, q database.Querier) (*response.PlatformStats, error)
}

type SubscriptionPlanRepository interface {
	FindBySlug(ctx context.Context, q database.Querier, slug string) (*entity.SubscriptionPlan, error)
	ListActive(ctx context.Context, q database.Querier) ([]entity.SubscriptionPlan, error)
}

type TenantSubscriptionRepository interface {
	Create(ctx context.Context, q database.Querier, s *entity.TenantSubscription) error
	FindActiveByTenant(ctx context.Context, q database.Querier, tenantID string) (*entity.TenantSubscription, error)
}

type SuperAdminRepository interface {
	FindByEmail(ctx context.Context, q database.Querier, email string) (*entity.SuperAdmin, error)
	UpdateLastLogin(ctx context.Context, q database.Querier, id string) error
}
