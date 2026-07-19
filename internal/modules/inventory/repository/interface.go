// Package repository defines the contracts the inventory service depends on.
package repository

import (
	"context"
	"errors"
	"time"

	"rentos/internal/modules/inventory/entity"
	"rentos/pkg/database"
)

// ErrNotFound is returned when a queried row does not exist.
var ErrNotFound = errors.New("repository: record not found")

// CategoryRepository manages the categories table.
type CategoryRepository interface {
	Create(ctx context.Context, q database.Querier, c *entity.Category) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Category, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Category, error)
	Update(ctx context.Context, q database.Querier, c *entity.Category) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// AssetTemplateRepository manages the asset_templates table.
type AssetTemplateRepository interface {
	Create(ctx context.Context, q database.Querier, t *entity.AssetTemplate) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.AssetTemplate, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.AssetTemplate, error)
	Update(ctx context.Context, q database.Querier, t *entity.AssetTemplate) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// TemplateFieldRepository manages the template_fields table.
type TemplateFieldRepository interface {
	Create(ctx context.Context, q database.Querier, f *entity.TemplateField) error
	ListByTemplate(ctx context.Context, q database.Querier, templateID, tenantID string) ([]entity.TemplateField, error)
}

// AssetRepository manages the assets table.
type AssetRepository interface {
	Create(ctx context.Context, q database.Querier, a *entity.Asset) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Asset, error)
	List(ctx context.Context, q database.Querier, tenantID string, categoryID, status, search *string, limit, offset int) ([]entity.Asset, error)
	Update(ctx context.Context, q database.Querier, a *entity.Asset) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// AssetValueRepository manages the asset_values table (EAV).
type AssetValueRepository interface {
	Upsert(ctx context.Context, q database.Querier, v *entity.AssetValue) error
	ListByAsset(ctx context.Context, q database.Querier, assetID, tenantID string) ([]entity.AssetValue, error)
}

// AssetImageRepository manages the asset_images table.
type AssetImageRepository interface {
	Create(ctx context.Context, q database.Querier, img *entity.AssetImage) error
	ListByAsset(ctx context.Context, q database.Querier, assetID, tenantID string) ([]entity.AssetImage, error)
	Delete(ctx context.Context, q database.Querier, id, assetID, actorID string) error
}

// AssetDocumentRepository manages the asset_documents table.
type AssetDocumentRepository interface {
	Create(ctx context.Context, q database.Querier, doc *entity.AssetDocument) error
	ListByAsset(ctx context.Context, q database.Querier, assetID, tenantID string) ([]entity.AssetDocument, error)
	Delete(ctx context.Context, q database.Querier, id, assetID, actorID string) error
}

// AssetAvailabilityRepository manages the asset_availability table.
type AssetAvailabilityRepository interface {
	Create(ctx context.Context, q database.Querier, a *entity.AssetAvailability) error
	CountConflicts(ctx context.Context, q database.Querier, assetID string, start, end time.Time) (int, error)
	ListByAsset(ctx context.Context, q database.Querier, assetID, tenantID string) ([]entity.AssetAvailability, error)
}
