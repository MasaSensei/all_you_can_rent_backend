package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos-backend/internal/modules/cms/entity"
	"rentos-backend/internal/modules/cms/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// websiteRepository
// ============================================================

type websiteRepository struct {
	qCreate   string
	qFindByID string
	qList     string
	qUpdate   string
	qDelete   string
}

func NewWebsiteRepository(qCreate, qFindByID, qList, qUpdate, qDelete string) repository.WebsiteRepository {
	return &websiteRepository{qCreate: qCreate, qFindByID: qFindByID, qList: qList, qUpdate: qUpdate, qDelete: qDelete}
}

func (r *websiteRepository) Create(ctx context.Context, q database.Querier, w *entity.Website) error {
	_, err := q.ExecContext(ctx, r.qCreate, w.ID, w.TenantID, w.Domain, w.Title, w.Theme, w.CreatedBy)
	if err != nil {
		return fmt.Errorf("websiteRepository.Create: %w", err)
	}
	return nil
}

func (r *websiteRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Website, error) {
	var w entity.Website
	if err := q.GetContext(ctx, &w, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("websiteRepository.FindByID: %w", err)
	}
	return &w, nil
}

func (r *websiteRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Website, error) {
	var out []entity.Website
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("websiteRepository.List: %w", err)
	}
	return out, nil
}

func (r *websiteRepository) Update(ctx context.Context, q database.Querier, w *entity.Website) error {
	res, err := q.ExecContext(ctx, r.qUpdate, w.ID, w.TenantID, w.Title, w.Theme)
	if err != nil {
		return fmt.Errorf("websiteRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *websiteRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("websiteRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// pageRepository
// ============================================================

type pageRepository struct {
	qCreate      string
	qFindByID    string
	qFindBySlug  string
	qList        string
	qUpdate      string
	qDelete      string
}

func NewPageRepository(qCreate, qFindByID, qFindBySlug, qList, qUpdate, qDelete string) repository.PageRepository {
	return &pageRepository{
		qCreate: qCreate, qFindByID: qFindByID, qFindBySlug: qFindBySlug,
		qList: qList, qUpdate: qUpdate, qDelete: qDelete,
	}
}

func (r *pageRepository) Create(ctx context.Context, q database.Querier, p *entity.Page) error {
	_, err := q.ExecContext(ctx, r.qCreate, p.ID, p.TenantID, p.WebsiteID, p.Title, p.Slug, p.Content, p.Template, p.CreatedBy)
	if err != nil {
		return fmt.Errorf("pageRepository.Create: %w", err)
	}
	return nil
}

func (r *pageRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Page, error) {
	var p entity.Page
	if err := q.GetContext(ctx, &p, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("pageRepository.FindByID: %w", err)
	}
	return &p, nil
}

func (r *pageRepository) FindBySlug(ctx context.Context, q database.Querier, websiteID, slug string) (*entity.Page, error) {
	var p entity.Page
	if err := q.GetContext(ctx, &p, r.qFindBySlug, websiteID, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("pageRepository.FindBySlug: %w", err)
	}
	return &p, nil
}

func (r *pageRepository) List(ctx context.Context, q database.Querier, websiteID, tenantID string, limit, offset int) ([]entity.Page, error) {
	var out []entity.Page
	if err := q.SelectContext(ctx, &out, r.qList, websiteID, tenantID, limit, offset); err != nil {
		return nil, fmt.Errorf("pageRepository.List: %w", err)
	}
	return out, nil
}

func (r *pageRepository) Update(ctx context.Context, q database.Querier, p *entity.Page) error {
	res, err := q.ExecContext(ctx, r.qUpdate, p.ID, p.TenantID, p.Title, p.Content, p.Template, p.Status)
	if err != nil {
		return fmt.Errorf("pageRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *pageRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("pageRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// menuRepository
// ============================================================

type menuRepository struct {
	qCreate   string
	qFindByID string
	qList     string
	qDelete   string
}

func NewMenuRepository(qCreate, qFindByID, qList, qDelete string) repository.MenuRepository {
	return &menuRepository{qCreate: qCreate, qFindByID: qFindByID, qList: qList, qDelete: qDelete}
}

func (r *menuRepository) Create(ctx context.Context, q database.Querier, m *entity.Menu) error {
	_, err := q.ExecContext(ctx, r.qCreate, m.ID, m.TenantID, m.WebsiteID, m.Name, m.Location, m.CreatedBy)
	if err != nil {
		return fmt.Errorf("menuRepository.Create: %w", err)
	}
	return nil
}

func (r *menuRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Menu, error) {
	var m entity.Menu
	if err := q.GetContext(ctx, &m, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("menuRepository.FindByID: %w", err)
	}
	return &m, nil
}

func (r *menuRepository) List(ctx context.Context, q database.Querier, websiteID, tenantID string) ([]entity.Menu, error) {
	var out []entity.Menu
	if err := q.SelectContext(ctx, &out, r.qList, websiteID, tenantID); err != nil {
		return nil, fmt.Errorf("menuRepository.List: %w", err)
	}
	return out, nil
}

func (r *menuRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("menuRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// menuItemRepository
// ============================================================

type menuItemRepository struct {
	qCreate      string
	qListByMenu  string
	qDelete      string
}

func NewMenuItemRepository(qCreate, qListByMenu, qDelete string) repository.MenuItemRepository {
	return &menuItemRepository{qCreate: qCreate, qListByMenu: qListByMenu, qDelete: qDelete}
}

func (r *menuItemRepository) Create(ctx context.Context, q database.Querier, item *entity.MenuItem) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		item.ID, item.TenantID, item.MenuID, item.ParentID,
		item.Label, item.URL, item.SortOrder, item.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("menuItemRepository.Create: %w", err)
	}
	return nil
}

func (r *menuItemRepository) ListByMenu(ctx context.Context, q database.Querier, menuID, tenantID string) ([]entity.MenuItem, error) {
	var out []entity.MenuItem
	if err := q.SelectContext(ctx, &out, r.qListByMenu, menuID, tenantID); err != nil {
		return nil, fmt.Errorf("menuItemRepository.ListByMenu: %w", err)
	}
	return out, nil
}

func (r *menuItemRepository) Delete(ctx context.Context, q database.Querier, id, menuID, actorID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, menuID, actorID)
	if err != nil {
		return fmt.Errorf("menuItemRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}
