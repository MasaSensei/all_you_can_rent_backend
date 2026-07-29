package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"rentos-backend/internal/modules/inventory/entity"
	"rentos-backend/internal/modules/inventory/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// assetRepository
// ============================================================

type assetRepository struct {
	qCreate   string
	qFindByID string
	qList     string
	qUpdate   string
	qDelete   string
}

func NewAssetRepository(qCreate, qFindByID, qList, qUpdate, qDelete string) repository.AssetRepository {
	return &assetRepository{
		qCreate: qCreate, qFindByID: qFindByID,
		qList: qList, qUpdate: qUpdate, qDelete: qDelete,
	}
}

func (r *assetRepository) Create(ctx context.Context, q database.Querier, a *entity.Asset) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		a.ID, a.TenantID, a.CategoryID, a.AssetTemplateID,
		a.Name, a.SKU, a.SerialNumber, a.Description,
		a.PurchasePrice, a.ReplacementValue, a.PurchaseDate,
		a.Condition, a.Location, a.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("assetRepository.Create: %w", err)
	}
	return nil
}

func (r *assetRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Asset, error) {
	var a entity.Asset
	if err := q.GetContext(ctx, &a, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("assetRepository.FindByID: %w", err)
	}
	return &a, nil
}

func (r *assetRepository) List(ctx context.Context, q database.Querier, tenantID string, categoryID, status, search *string, limit, offset int) ([]entity.Asset, error) {
	var out []entity.Asset
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, categoryID, status, search, limit, offset); err != nil {
		return nil, fmt.Errorf("assetRepository.List: %w", err)
	}
	return out, nil
}

func (r *assetRepository) Update(ctx context.Context, q database.Querier, a *entity.Asset) error {
	res, err := q.ExecContext(ctx, r.qUpdate,
		a.ID, a.TenantID, a.Name, a.Description, a.Condition,
		a.Location, a.PurchasePrice, a.ReplacementValue,
	)
	if err != nil {
		return fmt.Errorf("assetRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *assetRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("assetRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// assetValueRepository
// ============================================================

type assetValueRepository struct {
	qUpsert string
	qList   string
}

func NewAssetValueRepository(qUpsert, qList string) repository.AssetValueRepository {
	return &assetValueRepository{qUpsert: qUpsert, qList: qList}
}

func (r *assetValueRepository) Upsert(ctx context.Context, q database.Querier, v *entity.AssetValue) error {
	_, err := q.ExecContext(ctx, r.qUpsert,
		v.ID, v.TenantID, v.AssetID, v.TemplateFieldID, v.Value,
	)
	if err != nil {
		return fmt.Errorf("assetValueRepository.Upsert: %w", err)
	}
	return nil
}

func (r *assetValueRepository) ListByAsset(ctx context.Context, q database.Querier, assetID, tenantID string) ([]entity.AssetValue, error) {
	var out []entity.AssetValue
	if err := q.SelectContext(ctx, &out, r.qList, assetID, tenantID); err != nil {
		return nil, fmt.Errorf("assetValueRepository.ListByAsset: %w", err)
	}
	return out, nil
}

// ============================================================
// assetImageRepository
// ============================================================

type assetImageRepository struct {
	qCreate string
	qList   string
	qDelete string
}

func NewAssetImageRepository(qCreate, qList, qDelete string) repository.AssetImageRepository {
	return &assetImageRepository{qCreate: qCreate, qList: qList, qDelete: qDelete}
}

func (r *assetImageRepository) Create(ctx context.Context, q database.Querier, img *entity.AssetImage) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		img.ID, img.TenantID, img.AssetID, img.URL,
		img.AltText, img.IsPrimary, img.SortOrder, img.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("assetImageRepository.Create: %w", err)
	}
	return nil
}

func (r *assetImageRepository) ListByAsset(ctx context.Context, q database.Querier, assetID, tenantID string) ([]entity.AssetImage, error) {
	var out []entity.AssetImage
	if err := q.SelectContext(ctx, &out, r.qList, assetID, tenantID); err != nil {
		return nil, fmt.Errorf("assetImageRepository.ListByAsset: %w", err)
	}
	return out, nil
}

func (r *assetImageRepository) Delete(ctx context.Context, q database.Querier, id, assetID, actorID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, assetID, actorID)
	if err != nil {
		return fmt.Errorf("assetImageRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// assetDocumentRepository
// ============================================================

type assetDocumentRepository struct {
	qCreate string
	qList   string
	qDelete string
}

func NewAssetDocumentRepository(qCreate, qList, qDelete string) repository.AssetDocumentRepository {
	return &assetDocumentRepository{qCreate: qCreate, qList: qList, qDelete: qDelete}
}

func (r *assetDocumentRepository) Create(ctx context.Context, q database.Querier, doc *entity.AssetDocument) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		doc.ID, doc.TenantID, doc.AssetID, doc.Title,
		doc.FileURL, doc.FileType, doc.FileSize, doc.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("assetDocumentRepository.Create: %w", err)
	}
	return nil
}

func (r *assetDocumentRepository) ListByAsset(ctx context.Context, q database.Querier, assetID, tenantID string) ([]entity.AssetDocument, error) {
	var out []entity.AssetDocument
	if err := q.SelectContext(ctx, &out, r.qList, assetID, tenantID); err != nil {
		return nil, fmt.Errorf("assetDocumentRepository.ListByAsset: %w", err)
	}
	return out, nil
}

func (r *assetDocumentRepository) Delete(ctx context.Context, q database.Querier, id, assetID, actorID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, assetID, actorID)
	if err != nil {
		return fmt.Errorf("assetDocumentRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// assetAvailabilityRepository
// ============================================================

type assetAvailabilityRepository struct {
	qCreate         string
	qCountConflicts string
	qList           string
}

func NewAssetAvailabilityRepository(qCreate, qCountConflicts, qList string) repository.AssetAvailabilityRepository {
	return &assetAvailabilityRepository{
		qCreate: qCreate, qCountConflicts: qCountConflicts, qList: qList,
	}
}

func (r *assetAvailabilityRepository) Create(ctx context.Context, q database.Querier, a *entity.AssetAvailability) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		a.ID, a.TenantID, a.AssetID, a.StartDate, a.EndDate,
		a.AvailabilityType, a.Reason, a.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("assetAvailabilityRepository.Create: %w", err)
	}
	return nil
}

func (r *assetAvailabilityRepository) CountConflicts(ctx context.Context, q database.Querier, assetID string, start, end time.Time) (int, error) {
	var count int
	if err := q.GetContext(ctx, &count, r.qCountConflicts, assetID, start, end); err != nil {
		return 0, fmt.Errorf("assetAvailabilityRepository.CountConflicts: %w", err)
	}
	return count, nil
}

func (r *assetAvailabilityRepository) ListByAsset(ctx context.Context, q database.Querier, assetID, tenantID string) ([]entity.AssetAvailability, error) {
	var out []entity.AssetAvailability
	if err := q.SelectContext(ctx, &out, r.qList, assetID, tenantID); err != nil {
		return nil, fmt.Errorf("assetAvailabilityRepository.ListByAsset: %w", err)
	}
	return out, nil
}
