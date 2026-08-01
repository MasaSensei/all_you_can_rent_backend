package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/cms/dto/request"
	"rentos-backend/internal/modules/cms/dto/response"
	"rentos-backend/internal/modules/cms/entity"
	"rentos-backend/internal/modules/cms/repository"
	pkgresponse "rentos-backend/pkg/response"
)

// ============================================================
// websiteService
// ============================================================

type websiteService struct {
	db   *sqlx.DB
	repo repository.WebsiteRepository
}

func NewWebsiteService(db *sqlx.DB, repo repository.WebsiteRepository) WebsiteService {
	return &websiteService{db: db, repo: repo}
}

func (s *websiteService) Create(ctx context.Context, tenantID, actorID string, req request.CreateWebsite) (*response.Website, error) {
	w := &entity.Website{
		ID: uuid.NewString(), TenantID: tenantID,
		Domain: req.Domain, Title: req.Title, Theme: req.Theme, CreatedBy: &actorID,
	}
	if err := s.repo.Create(ctx, s.db, w); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, w.ID, tenantID)
}

func (s *websiteService) GetByID(ctx context.Context, id, tenantID string) (*response.Website, error) {
	w, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "website not found")
		}
		return nil, err
	}
	return toWebsiteResponse(w), nil
}

func (s *websiteService) List(ctx context.Context, tenantID string) ([]response.Website, error) {
	websites, err := s.repo.List(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.Website, 0, len(websites))
	for _, w := range websites {
		out = append(out, *toWebsiteResponse(&w))
	}
	return out, nil
}

func (s *websiteService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateWebsite) (*response.Website, error) {
	w := &entity.Website{ID: id, TenantID: tenantID, Title: derefStr(req.Title), Theme: req.Theme, UpdatedBy: &actorID}
	if err := s.repo.Update(ctx, s.db, w); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "website not found")
		}
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

func (s *websiteService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.repo.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "website not found")
		}
		return err
	}
	return nil
}

// ============================================================
// pageService
// ============================================================

type pageService struct {
	db   *sqlx.DB
	repo repository.PageRepository
}

func NewPageService(db *sqlx.DB, repo repository.PageRepository) PageService {
	return &pageService{db: db, repo: repo}
}

func (s *pageService) Create(ctx context.Context, tenantID, actorID string, req request.CreatePage) (*response.Page, error) {
	p := &entity.Page{
		ID: uuid.NewString(), TenantID: tenantID, WebsiteID: req.WebsiteID,
		Title: req.Title, Slug: req.Slug, Content: req.Content, Template: req.Template, CreatedBy: &actorID,
	}
	if err := s.repo.Create(ctx, s.db, p); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, p.ID, tenantID)
}

func (s *pageService) GetByID(ctx context.Context, id, tenantID string) (*response.Page, error) {
	p, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "page not found")
		}
		return nil, err
	}
	return toPageResponse(p), nil
}

func (s *pageService) GetBySlug(ctx context.Context, websiteID, slug string) (*response.Page, error) {
	p, err := s.repo.FindBySlug(ctx, s.db, websiteID, slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "page not found")
		}
		return nil, err
	}
	return toPageResponse(p), nil
}

func (s *pageService) List(ctx context.Context, websiteID, tenantID string, page, perPage int) ([]response.Page, error) {
	perPage, page = normPage(perPage, page)
	pages, err := s.repo.List(ctx, s.db, websiteID, tenantID, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.Page, 0, len(pages))
	for _, p := range pages {
		out = append(out, *toPageResponse(&p))
	}
	return out, nil
}

func (s *pageService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdatePage) (*response.Page, error) {
	p := &entity.Page{ID: id, TenantID: tenantID, Title: derefStr(req.Title), Content: req.Content, Template: req.Template, Status: derefStr(req.Status), UpdatedBy: &actorID}
	if err := s.repo.Update(ctx, s.db, p); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "page not found")
		}
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

func (s *pageService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.repo.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "page not found")
		}
		return err
	}
	return nil
}

// ============================================================
// menuService
// ============================================================

type menuService struct {
	db        *sqlx.DB
	menus     repository.MenuRepository
	menuItems repository.MenuItemRepository
}

func NewMenuService(db *sqlx.DB, menus repository.MenuRepository, items repository.MenuItemRepository) MenuService {
	return &menuService{db: db, menus: menus, menuItems: items}
}

func (s *menuService) Create(ctx context.Context, tenantID, actorID string, req request.CreateMenu) (*response.Menu, error) {
	m := &entity.Menu{
		ID: uuid.NewString(), TenantID: tenantID, WebsiteID: req.WebsiteID,
		Name: req.Name, Location: req.Location, CreatedBy: &actorID,
	}
	if err := s.menus.Create(ctx, s.db, m); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, m.ID, tenantID)
}

func (s *menuService) GetByID(ctx context.Context, id, tenantID string) (*response.Menu, error) {
	m, err := s.menus.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "menu not found")
		}
		return nil, err
	}
	items, _ := s.menuItems.ListByMenu(ctx, s.db, m.ID, tenantID)
	return toMenuResponse(m, buildItemTree(items)), nil
}

func (s *menuService) List(ctx context.Context, websiteID, tenantID string) ([]response.Menu, error) {
	menus, err := s.menus.List(ctx, s.db, websiteID, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.Menu, 0, len(menus))
	for _, m := range menus {
		items, _ := s.menuItems.ListByMenu(ctx, s.db, m.ID, tenantID)
		out = append(out, *toMenuResponse(&m, buildItemTree(items)))
	}
	return out, nil
}

func (s *menuService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.menus.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "menu not found")
		}
		return err
	}
	return nil
}

func (s *menuService) AddItem(ctx context.Context, menuID, tenantID, actorID string, req request.AddMenuItem) (*response.MenuItem, error) {
	item := &entity.MenuItem{
		ID: uuid.NewString(), TenantID: tenantID, MenuID: menuID,
		ParentID: req.ParentID, Label: req.Label, URL: req.URL, SortOrder: req.SortOrder, CreatedBy: &actorID,
	}
	if err := s.menuItems.Create(ctx, s.db, item); err != nil {
		return nil, err
	}
	return &response.MenuItem{ID: item.ID, ParentID: item.ParentID, Label: item.Label, URL: item.URL, SortOrder: item.SortOrder}, nil
}

func (s *menuService) DeleteItem(ctx context.Context, itemID, menuID, tenantID, actorID string) error {
	if err := s.menuItems.Delete(ctx, s.db, itemID, menuID, actorID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "menu item not found")
		}
		return err
	}
	return nil
}

// ============================================================
// blogService
// ============================================================

type blogService struct {
	db         *sqlx.DB
	blogs      repository.BlogRepository
	categories repository.BlogCategoryRepository
}

func NewBlogService(db *sqlx.DB, blogs repository.BlogRepository, categories repository.BlogCategoryRepository) BlogService {
	return &blogService{db: db, blogs: blogs, categories: categories}
}

func (s *blogService) CreateCategory(ctx context.Context, tenantID, actorID string, req request.CreateBlogCategory) (*response.BlogCategory, error) {
	c := &entity.BlogCategory{ID: uuid.NewString(), TenantID: tenantID, Name: req.Name, Slug: req.Slug, CreatedBy: &actorID}
	if err := s.categories.Create(ctx, s.db, c); err != nil {
		return nil, err
	}
	return &response.BlogCategory{ID: c.ID, TenantID: c.TenantID, Name: c.Name, Slug: c.Slug, Status: "active", CreatedAt: c.CreatedAt}, nil
}

func (s *blogService) ListCategories(ctx context.Context, tenantID string) ([]response.BlogCategory, error) {
	cats, err := s.categories.List(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.BlogCategory, 0, len(cats))
	for _, c := range cats {
		out = append(out, response.BlogCategory{ID: c.ID, TenantID: c.TenantID, Name: c.Name, Slug: c.Slug, Status: c.Status, CreatedAt: c.CreatedAt})
	}
	return out, nil
}

func (s *blogService) DeleteCategory(ctx context.Context, id, tenantID string) error {
	if err := s.categories.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "blog category not found")
		}
		return err
	}
	return nil
}

func (s *blogService) Create(ctx context.Context, tenantID, actorID string, req request.CreateBlog) (*response.Blog, error) {
	b := &entity.Blog{
		ID: uuid.NewString(), TenantID: tenantID, WebsiteID: req.WebsiteID,
		AuthorID: &actorID, BlogCategoryID: req.BlogCategoryID,
		Title: req.Title, Slug: req.Slug, Content: req.Content, FeaturedImage: req.FeaturedImage, CreatedBy: &actorID,
	}
	if err := s.blogs.Create(ctx, s.db, b); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, b.ID, tenantID)
}

func (s *blogService) GetByID(ctx context.Context, id, tenantID string) (*response.Blog, error) {
	b, err := s.blogs.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "blog not found")
		}
		return nil, err
	}
	return toBlogResponse(b), nil
}

func (s *blogService) GetBySlug(ctx context.Context, websiteID, slug string) (*response.Blog, error) {
	b, err := s.blogs.FindBySlug(ctx, s.db, websiteID, slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "blog not found")
		}
		return nil, err
	}
	return toBlogResponse(b), nil
}

func (s *blogService) List(ctx context.Context, tenantID string, filter request.ListBlogsFilter) ([]response.Blog, error) {
	perPage, page := normPage(filter.PerPage, filter.Page)
	blogs, err := s.blogs.List(ctx, s.db, tenantID, filter.WebsiteID, filter.BlogCategoryID, filter.Status, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.Blog, 0, len(blogs))
	for _, b := range blogs {
		out = append(out, *toBlogResponse(&b))
	}
	return out, nil
}

func (s *blogService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateBlog) (*response.Blog, error) {
	b := &entity.Blog{ID: id, TenantID: tenantID, Title: derefStr(req.Title), Content: req.Content, FeaturedImage: req.FeaturedImage, Status: derefStr(req.Status), UpdatedBy: &actorID}
	if err := s.blogs.Update(ctx, s.db, b); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "blog not found")
		}
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

func (s *blogService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.blogs.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "blog not found")
		}
		return err
	}
	return nil
}

// ============================================================
// seoService
// ============================================================

type seoService struct {
	db   *sqlx.DB
	repo repository.SEOMetaRepository
}

func NewSEOService(db *sqlx.DB, repo repository.SEOMetaRepository) SEOService {
	return &seoService{db: db, repo: repo}
}

func (s *seoService) Upsert(ctx context.Context, tenantID string, req request.UpsertSEO) (*response.SEOMeta, error) {
	meta := &entity.SEOMeta{
		ID: uuid.NewString(), TenantID: tenantID,
		EntityType: req.EntityType, EntityID: req.EntityID,
		MetaTitle: req.MetaTitle, MetaDescription: req.MetaDescription,
		MetaKeywords: req.MetaKeywords, OGImage: req.OGImage, CanonicalURL: req.CanonicalURL,
	}
	if err := s.repo.Upsert(ctx, s.db, meta); err != nil {
		return nil, err
	}
	return s.Get(ctx, req.EntityType, req.EntityID)
}

func (s *seoService) Get(ctx context.Context, entityType, entityID string) (*response.SEOMeta, error) {
	meta, err := s.repo.Find(ctx, s.db, entityType, entityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "SEO meta not found")
		}
		return nil, err
	}
	return &response.SEOMeta{
		ID: meta.ID, EntityType: meta.EntityType, EntityID: meta.EntityID,
		MetaTitle: meta.MetaTitle, MetaDescription: meta.MetaDescription,
		MetaKeywords: meta.MetaKeywords, OGImage: meta.OGImage, CanonicalURL: meta.CanonicalURL,
	}, nil
}

// ============================================================
// helpers
// ============================================================

func buildItemTree(items []entity.MenuItem) []response.MenuItem {
	index := make(map[string]*response.MenuItem, len(items))
	for _, item := range items {
		cp := response.MenuItem{ID: item.ID, ParentID: item.ParentID, Label: item.Label, URL: item.URL, SortOrder: item.SortOrder}
		index[item.ID] = &cp
	}
	var roots []response.MenuItem
	for _, item := range items {
		node := index[item.ID]
		if item.ParentID == nil {
			roots = append(roots, *node)
		} else if parent, ok := index[*item.ParentID]; ok {
			parent.Children = append(parent.Children, *node)
		}
	}
	return roots
}

func toWebsiteResponse(w *entity.Website) *response.Website {
	return &response.Website{ID: w.ID, TenantID: w.TenantID, Domain: w.Domain, Title: w.Title, Theme: w.Theme, Status: w.Status, CreatedAt: w.CreatedAt, UpdatedAt: w.UpdatedAt}
}

func toPageResponse(p *entity.Page) *response.Page {
	return &response.Page{ID: p.ID, TenantID: p.TenantID, WebsiteID: p.WebsiteID, Title: p.Title, Slug: p.Slug, Content: p.Content, Template: p.Template, Status: p.Status, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

func toMenuResponse(m *entity.Menu, items []response.MenuItem) *response.Menu {
	return &response.Menu{ID: m.ID, TenantID: m.TenantID, WebsiteID: m.WebsiteID, Name: m.Name, Location: m.Location, Status: m.Status, Items: items, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func toBlogResponse(b *entity.Blog) *response.Blog {
	return &response.Blog{ID: b.ID, TenantID: b.TenantID, WebsiteID: b.WebsiteID, AuthorID: b.AuthorID, BlogCategoryID: b.BlogCategoryID, Title: b.Title, Slug: b.Slug, Content: b.Content, FeaturedImage: b.FeaturedImage, PublishedAt: b.PublishedAt, Status: b.Status, CreatedAt: b.CreatedAt, UpdatedAt: b.UpdatedAt}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
