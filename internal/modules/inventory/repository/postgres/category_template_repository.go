package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos/internal/modules/inventory/entity"
	"rentos/internal/modules/inventory/repository"
	"rentos/pkg/database"
)

// ============================================================
// categoryRepository
// ============================================================

type categoryRepository struct {
	qCreate   string
	qFindByID string
	qList     string
	qUpdate   string
	qDelete   string
}

func NewCategoryRepository(qCreate, qFindByID, qList, qUpdate, qDelete string) repository.CategoryRepository {
	return &categoryRepository{
		qCreate: qCreate, qFindByID: qFindByID,
		qList: qList, qUpdate: qUpdate, qDelete: qDelete,
	}
}

func (r *categoryRepository) Create(ctx context.Context, q database.Querier, c *entity.Category) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		c.ID, c.TenantID, c.ParentID, c.Name, c.Slug,
		c.Description, c.Icon, c.SortOrder, c.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("categoryRepository.Create: %w", err)
	}
	return nil
}

func (r *categoryRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Category, error) {
	var c entity.Category
	if err := q.GetContext(ctx, &c, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("categoryRepository.FindByID: %w", err)
	}
	return &c, nil
}

func (r *categoryRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Category, error) {
	var out []entity.Category
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("categoryRepository.List: %w", err)
	}
	return out, nil
}

func (r *categoryRepository) Update(ctx context.Context, q database.Querier, c *entity.Category) error {
	res, err := q.ExecContext(ctx, r.qUpdate,
		c.ID, c.TenantID, c.Name, c.Description, c.Icon, c.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("categoryRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("categoryRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// assetTemplateRepository
// ============================================================

type assetTemplateRepository struct {
	qCreate   string
	qFindByID string
	qList     string
	qUpdate   string
	qDelete   string
}

func NewAssetTemplateRepository(qCreate, qFindByID, qList, qUpdate, qDelete string) repository.AssetTemplateRepository {
	return &assetTemplateRepository{
		qCreate: qCreate, qFindByID: qFindByID,
		qList: qList, qUpdate: qUpdate, qDelete: qDelete,
	}
}

func (r *assetTemplateRepository) Create(ctx context.Context, q database.Querier, t *entity.AssetTemplate) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		t.ID, t.TenantID, t.CategoryID, t.Name, t.Description, t.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("assetTemplateRepository.Create: %w", err)
	}
	return nil
}

func (r *assetTemplateRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.AssetTemplate, error) {
	var t entity.AssetTemplate
	if err := q.GetContext(ctx, &t, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("assetTemplateRepository.FindByID: %w", err)
	}
	return &t, nil
}

func (r *assetTemplateRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.AssetTemplate, error) {
	var out []entity.AssetTemplate
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("assetTemplateRepository.List: %w", err)
	}
	return out, nil
}

func (r *assetTemplateRepository) Update(ctx context.Context, q database.Querier, t *entity.AssetTemplate) error {
	res, err := q.ExecContext(ctx, r.qUpdate,
		t.ID, t.TenantID, t.Name, t.Description,
	)
	if err != nil {
		return fmt.Errorf("assetTemplateRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *assetTemplateRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("assetTemplateRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// templateFieldRepository
// ============================================================

type templateFieldRepository struct {
	qCreate string
	qList   string
}

func NewTemplateFieldRepository(qCreate, qList string) repository.TemplateFieldRepository {
	return &templateFieldRepository{qCreate: qCreate, qList: qList}
}

func (r *templateFieldRepository) Create(ctx context.Context, q database.Querier, f *entity.TemplateField) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		f.ID, f.TenantID, f.AssetTemplateID, f.FieldName, f.FieldLabel,
		f.FieldType, f.IsRequired, f.DefaultValue, f.Options, f.SortOrder, f.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("templateFieldRepository.Create: %w", err)
	}
	return nil
}

func (r *templateFieldRepository) ListByTemplate(ctx context.Context, q database.Querier, templateID, tenantID string) ([]entity.TemplateField, error) {
	var out []entity.TemplateField
	if err := q.SelectContext(ctx, &out, r.qList, templateID, tenantID); err != nil {
		return nil, fmt.Errorf("templateFieldRepository.ListByTemplate: %w", err)
	}
	return out, nil
}
