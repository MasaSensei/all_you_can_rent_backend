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
// blogCategoryRepository
// ============================================================

type blogCategoryRepository struct {
	qCreate string
	qList   string
	qDelete string
}

func NewBlogCategoryRepository(qCreate, qList, qDelete string) repository.BlogCategoryRepository {
	return &blogCategoryRepository{qCreate: qCreate, qList: qList, qDelete: qDelete}
}

func (r *blogCategoryRepository) Create(ctx context.Context, q database.Querier, c *entity.BlogCategory) error {
	_, err := q.ExecContext(ctx, r.qCreate, c.ID, c.TenantID, c.Name, c.Slug, c.CreatedBy)
	if err != nil {
		return fmt.Errorf("blogCategoryRepository.Create: %w", err)
	}
	return nil
}

func (r *blogCategoryRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.BlogCategory, error) {
	var out []entity.BlogCategory
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("blogCategoryRepository.List: %w", err)
	}
	return out, nil
}

func (r *blogCategoryRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("blogCategoryRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// blogRepository
// ============================================================

type blogRepository struct {
	qCreate      string
	qFindByID    string
	qFindBySlug  string
	qList        string
	qUpdate      string
	qDelete      string
}

func NewBlogRepository(qCreate, qFindByID, qFindBySlug, qList, qUpdate, qDelete string) repository.BlogRepository {
	return &blogRepository{
		qCreate: qCreate, qFindByID: qFindByID, qFindBySlug: qFindBySlug,
		qList: qList, qUpdate: qUpdate, qDelete: qDelete,
	}
}

func (r *blogRepository) Create(ctx context.Context, q database.Querier, b *entity.Blog) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		b.ID, b.TenantID, b.WebsiteID, b.AuthorID, b.BlogCategoryID,
		b.Title, b.Slug, b.Content, b.FeaturedImage, b.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("blogRepository.Create: %w", err)
	}
	return nil
}

func (r *blogRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Blog, error) {
	var b entity.Blog
	if err := q.GetContext(ctx, &b, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("blogRepository.FindByID: %w", err)
	}
	return &b, nil
}

func (r *blogRepository) FindBySlug(ctx context.Context, q database.Querier, websiteID, slug string) (*entity.Blog, error) {
	var b entity.Blog
	if err := q.GetContext(ctx, &b, r.qFindBySlug, websiteID, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("blogRepository.FindBySlug: %w", err)
	}
	return &b, nil
}

func (r *blogRepository) List(ctx context.Context, q database.Querier, tenantID string, websiteID, categoryID, status *string, limit, offset int) ([]entity.Blog, error) {
	var out []entity.Blog
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, websiteID, categoryID, status, limit, offset); err != nil {
		return nil, fmt.Errorf("blogRepository.List: %w", err)
	}
	return out, nil
}

func (r *blogRepository) Update(ctx context.Context, q database.Querier, b *entity.Blog) error {
	res, err := q.ExecContext(ctx, r.qUpdate,
		b.ID, b.TenantID, b.Title, b.Content, b.FeaturedImage, b.Status,
	)
	if err != nil {
		return fmt.Errorf("blogRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *blogRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("blogRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// seoMetaRepository
// ============================================================

type seoMetaRepository struct {
	qUpsert string
	qFind   string
}

func NewSEOMetaRepository(qUpsert, qFind string) repository.SEOMetaRepository {
	return &seoMetaRepository{qUpsert: qUpsert, qFind: qFind}
}

func (r *seoMetaRepository) Upsert(ctx context.Context, q database.Querier, s *entity.SEOMeta) error {
	_, err := q.ExecContext(ctx, r.qUpsert,
		s.ID, s.TenantID, s.EntityType, s.EntityID,
		s.MetaTitle, s.MetaDescription, s.MetaKeywords, s.OGImage, s.CanonicalURL,
	)
	if err != nil {
		return fmt.Errorf("seoMetaRepository.Upsert: %w", err)
	}
	return nil
}

func (r *seoMetaRepository) Find(ctx context.Context, q database.Querier, entityType, entityID string) (*entity.SEOMeta, error) {
	var s entity.SEOMeta
	if err := q.GetContext(ctx, &s, r.qFind, entityType, entityID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("seoMetaRepository.Find: %w", err)
	}
	return &s, nil
}
