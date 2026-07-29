// Package repository defines the contracts the core service depends on.
// Concrete implementations live in repository/postgres. Services depend
// on these interfaces, never on the postgres package directly, so the
// data-access layer can be swapped or mocked in tests without touching
// business logic.
package repository

import (
	"context"
	"errors"

	"rentos-backend/internal/modules/core/entity"
	"rentos-backend/pkg/database"
)

// ErrNotFound is returned by any repository implementation when a queried
// row does not exist. Defined here at the interface level so services only
// ever import this package — never a concrete implementation such as
// repository/postgres — keeping the dependency arrow pointing inward.
var ErrNotFound = errors.New("repository: record not found")

type TenantRepository interface {
	Create(ctx context.Context, q database.Querier, t *entity.Tenant) error
	FindByID(ctx context.Context, q database.Querier, id string) (*entity.Tenant, error)
	FindBySlug(ctx context.Context, q database.Querier, slug string) (*entity.Tenant, error)
	UpdateStatus(ctx context.Context, q database.Querier, id, status string) error
}

type SettingRepository interface {
	Upsert(ctx context.Context, q database.Querier, s *entity.Setting) error
	Find(ctx context.Context, q database.Querier, tenantID, key string) (*entity.Setting, error)
	ListByTenant(ctx context.Context, q database.Querier, tenantID string) ([]entity.Setting, error)
}

type SystemSettingRepository interface {
	Upsert(ctx context.Context, q database.Querier, s *entity.SystemSetting) error
	Find(ctx context.Context, q database.Querier, key string) (*entity.SystemSetting, error)
}

type AuditLogRepository interface {
	Create(ctx context.Context, q database.Querier, a *entity.AuditLog) error
	List(ctx context.Context, q database.Querier, tenantID string, entityType, entityID, action *string, limit, offset int) ([]entity.AuditLog, error)
}
