// Package service holds the core module's business logic. Handlers depend
// on these interfaces, never on the concrete struct, so they can be unit
// tested with a mock implementation.
package service

import (
	"context"

	"rentos-backend/internal/modules/core/dto/request"
	"rentos-backend/internal/modules/core/entity"
)

// TenantService manages tenant lifecycle.
type TenantService interface {
	CreateTenant(ctx context.Context, req request.CreateTenant) (*entity.Tenant, error)
	GetTenantByID(ctx context.Context, id string) (*entity.Tenant, error)
	GetTenantBySlug(ctx context.Context, slug string) (*entity.Tenant, error)
}

// SettingService manages per-tenant settings and global system settings.
type SettingService interface {
	UpsertSetting(ctx context.Context, tenantID string, req request.UpsertSetting, actorID *string) (*entity.Setting, error)
	GetSetting(ctx context.Context, tenantID, key string) (*entity.Setting, error)
	ListSettings(ctx context.Context, tenantID string) ([]entity.Setting, error)

	UpsertSystemSetting(ctx context.Context, req request.UpsertSystemSetting) (*entity.SystemSetting, error)
	GetSystemSetting(ctx context.Context, key string) (*entity.SystemSetting, error)
}

// AuditService writes and reads the audit trail. Other modules' services
// depend on this interface to record audit entries as a side effect of
// their own mutations.
type AuditService interface {
	Write(ctx context.Context, tenantID string, req request.CreateAuditLog) error
	List(ctx context.Context, tenantID string, filter request.ListAuditLogsFilter) ([]entity.AuditLog, error)
}
