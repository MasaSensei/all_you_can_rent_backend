// Package service holds the inventory module's business logic. The
// CheckAvailability method on AssetService is the critical cross-module
// contract consumed by the booking module — it depends on this interface,
// never on the inventory postgres package directly.
package service

import (
	"context"
	"time"

	"rentos/internal/modules/inventory/dto/request"
	"rentos/internal/modules/inventory/dto/response"
	"rentos/internal/modules/inventory/entity"
)

// CategoryService manages the category hierarchy.
type CategoryService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateCategory) (*entity.Category, error)
	GetByID(ctx context.Context, id, tenantID string) (*entity.Category, error)
	// List returns a flat list; the handler/consumer builds the tree.
	List(ctx context.Context, tenantID string) ([]entity.Category, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateCategory) (*entity.Category, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// AssetTemplateService manages asset templates and their fields.
type AssetTemplateService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateAssetTemplate) (*response.AssetTemplate, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.AssetTemplate, error)
	List(ctx context.Context, tenantID string) ([]response.AssetTemplate, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateAssetTemplate) (*response.AssetTemplate, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// AssetService manages the rentable asset catalogue.
type AssetService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateAsset) (*response.Asset, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Asset, error)
	List(ctx context.Context, tenantID string, filter request.ListAssetsFilter) ([]response.Asset, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateAsset) (*response.Asset, error)
	Delete(ctx context.Context, id, tenantID string) error

	AddImage(ctx context.Context, tenantID, assetID, actorID string, img *entity.AssetImage) (*entity.AssetImage, error)
	DeleteImage(ctx context.Context, tenantID, assetID, imageID, actorID string) error

	AddDocument(ctx context.Context, tenantID, assetID, actorID string, doc *entity.AssetDocument) (*entity.AssetDocument, error)
	DeleteDocument(ctx context.Context, tenantID, assetID, docID, actorID string) error

	BlockAvailability(ctx context.Context, tenantID, assetID, actorID string, req request.BlockAvailability) (*entity.AssetAvailability, error)
	ListAvailability(ctx context.Context, assetID, tenantID string) ([]entity.AssetAvailability, error)

	// CheckAvailability is the cross-module contract consumed by booking.
	// Returns true when the asset has no blocking conflicts in the range.
	CheckAvailability(ctx context.Context, assetID string, start, end time.Time) (bool, error)
}
