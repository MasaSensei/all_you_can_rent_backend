package repository

import (
	"context"
	"errors"

	"rentos-backend/internal/modules/cms/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

type WebsiteRepository interface {
	Create(ctx context.Context, q database.Querier, w *entity.Website) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Website, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Website, error)
	Update(ctx context.Context, q database.Querier, w *entity.Website) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

type PageRepository interface {
	Create(ctx context.Context, q database.Querier, p *entity.Page) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Page, error)
	FindBySlug(ctx context.Context, q database.Querier, websiteID, slug string) (*entity.Page, error)
	List(ctx context.Context, q database.Querier, websiteID, tenantID string, limit, offset int) ([]entity.Page, error)
	Update(ctx context.Context, q database.Querier, p *entity.Page) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

type MenuRepository interface {
	Create(ctx context.Context, q database.Querier, m *entity.Menu) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Menu, error)
	List(ctx context.Context, q database.Querier, websiteID, tenantID string) ([]entity.Menu, error)
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

type MenuItemRepository interface {
	Create(ctx context.Context, q database.Querier, item *entity.MenuItem) error
	ListByMenu(ctx context.Context, q database.Querier, menuID, tenantID string) ([]entity.MenuItem, error)
	Delete(ctx context.Context, q database.Querier, id, menuID, actorID string) error
}

type BlogCategoryRepository interface {
	Create(ctx context.Context, q database.Querier, c *entity.BlogCategory) error
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.BlogCategory, error)
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

type BlogRepository interface {
	Create(ctx context.Context, q database.Querier, b *entity.Blog) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Blog, error)
	FindBySlug(ctx context.Context, q database.Querier, websiteID, slug string) (*entity.Blog, error)
	List(ctx context.Context, q database.Querier, tenantID string, websiteID, categoryID, status *string, limit, offset int) ([]entity.Blog, error)
	Update(ctx context.Context, q database.Querier, b *entity.Blog) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

type SEOMetaRepository interface {
	Upsert(ctx context.Context, q database.Querier, s *entity.SEOMeta) error
	Find(ctx context.Context, q database.Querier, entityType, entityID string) (*entity.SEOMeta, error)
}
