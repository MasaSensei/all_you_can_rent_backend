package service

import (
	"context"

	"rentos-backend/internal/modules/cms/dto/request"
	"rentos-backend/internal/modules/cms/dto/response"
)

// WebsiteService manages websites.
type WebsiteService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateWebsite) (*response.Website, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Website, error)
	List(ctx context.Context, tenantID string) ([]response.Website, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateWebsite) (*response.Website, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// PageService manages pages within a website.
type PageService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreatePage) (*response.Page, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Page, error)
	GetBySlug(ctx context.Context, websiteID, slug string) (*response.Page, error)
	List(ctx context.Context, websiteID, tenantID string, page, perPage int) ([]response.Page, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdatePage) (*response.Page, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// MenuService manages menus and their item trees.
type MenuService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateMenu) (*response.Menu, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Menu, error)
	List(ctx context.Context, websiteID, tenantID string) ([]response.Menu, error)
	Delete(ctx context.Context, id, tenantID string) error
	AddItem(ctx context.Context, menuID, tenantID, actorID string, req request.AddMenuItem) (*response.MenuItem, error)
	DeleteItem(ctx context.Context, itemID, menuID, tenantID, actorID string) error
}

// BlogService manages blog posts and their categories.
type BlogService interface {
	CreateCategory(ctx context.Context, tenantID, actorID string, req request.CreateBlogCategory) (*response.BlogCategory, error)
	ListCategories(ctx context.Context, tenantID string) ([]response.BlogCategory, error)
	DeleteCategory(ctx context.Context, id, tenantID string) error

	Create(ctx context.Context, tenantID, actorID string, req request.CreateBlog) (*response.Blog, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Blog, error)
	GetBySlug(ctx context.Context, websiteID, slug string) (*response.Blog, error)
	List(ctx context.Context, tenantID string, filter request.ListBlogsFilter) ([]response.Blog, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateBlog) (*response.Blog, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// SEOService manages SEO metadata for any entity.
type SEOService interface {
	Upsert(ctx context.Context, tenantID string, req request.UpsertSEO) (*response.SEOMeta, error)
	Get(ctx context.Context, entityType, entityID string) (*response.SEOMeta, error)
}
