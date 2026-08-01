package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/inventory/dto/request"
	"rentos-backend/internal/modules/inventory/dto/response"
	"rentos-backend/internal/modules/inventory/entity"
	"rentos-backend/internal/modules/inventory/repository"
	pkgresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/transaction"
)

// ============================================================
// categoryService
// ============================================================

type categoryService struct {
	db   *sqlx.DB
	repo repository.CategoryRepository
}

func NewCategoryService(db *sqlx.DB, repo repository.CategoryRepository) CategoryService {
	return &categoryService{db: db, repo: repo}
}

func (s *categoryService) Create(ctx context.Context, tenantID, actorID string, req request.CreateCategory) (*entity.Category, error) {
	c := &entity.Category{
		ID:          uuid.NewString(),
		TenantID:    tenantID,
		ParentID:    req.ParentID,
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
		CreatedBy:   &actorID,
	}
	if err := s.repo.Create(ctx, s.db, c); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, s.db, c.ID, tenantID)
}

func (s *categoryService) GetByID(ctx context.Context, id, tenantID string) (*entity.Category, error) {
	c, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "category not found")
		}
		return nil, err
	}
	return c, nil
}

func (s *categoryService) List(ctx context.Context, tenantID string) ([]entity.Category, error) {
	return s.repo.List(ctx, s.db, tenantID)
}

func (s *categoryService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateCategory) (*entity.Category, error) {
	c := &entity.Category{
		ID:          id,
		TenantID:    tenantID,
		Name:        *req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		SortOrder:   derefInt(req.SortOrder),
		UpdatedBy:   &actorID,
	}
	if err := s.repo.Update(ctx, s.db, c); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "category not found")
		}
		return nil, err
	}
	return s.repo.FindByID(ctx, s.db, id, tenantID)
}

func (s *categoryService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.repo.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "category not found")
		}
		return err
	}
	return nil
}

// ============================================================
// assetTemplateService
// ============================================================

type assetTemplateService struct {
	db        *sqlx.DB
	templates repository.AssetTemplateRepository
	fields    repository.TemplateFieldRepository
}

func NewAssetTemplateService(db *sqlx.DB, templates repository.AssetTemplateRepository, fields repository.TemplateFieldRepository) AssetTemplateService {
	return &assetTemplateService{db: db, templates: templates, fields: fields}
}

func (s *assetTemplateService) Create(ctx context.Context, tenantID, actorID string, req request.CreateAssetTemplate) (*response.AssetTemplate, error) {
	t := &entity.AssetTemplate{
		ID:          uuid.NewString(),
		TenantID:    tenantID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actorID,
	}

	if err := transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if err := s.templates.Create(ctx, tx, t); err != nil {
			return err
		}
		for i, f := range req.Fields {
			field := &entity.TemplateField{
				ID:              uuid.NewString(),
				TenantID:        tenantID,
				AssetTemplateID: t.ID,
				FieldName:       f.FieldName,
				FieldLabel:      f.FieldLabel,
				FieldType:       f.FieldType,
				IsRequired:      f.IsRequired,
				DefaultValue:    f.DefaultValue,
				Options:         f.Options,
				SortOrder:       i,
				CreatedBy:       &actorID,
			}
			if err := s.fields.Create(ctx, tx, field); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, t.ID, tenantID)
}

func (s *assetTemplateService) GetByID(ctx context.Context, id, tenantID string) (*response.AssetTemplate, error) {
	t, err := s.templates.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "asset template not found")
		}
		return nil, err
	}
	fields, err := s.fields.ListByTemplate(ctx, s.db, id, tenantID)
	if err != nil {
		return nil, err
	}
	return toTemplateResponse(t, fields), nil
}

func (s *assetTemplateService) List(ctx context.Context, tenantID string) ([]response.AssetTemplate, error) {
	templates, err := s.templates.List(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.AssetTemplate, 0, len(templates))
	for _, t := range templates {
		fields, _ := s.fields.ListByTemplate(ctx, s.db, t.ID, tenantID)
		out = append(out, *toTemplateResponse(&t, fields))
	}
	return out, nil
}

func (s *assetTemplateService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateAssetTemplate) (*response.AssetTemplate, error) {
	t := &entity.AssetTemplate{ID: id, TenantID: tenantID, Name: *req.Name, Description: req.Description, UpdatedBy: &actorID}
	if err := s.templates.Update(ctx, s.db, t); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "asset template not found")
		}
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

func (s *assetTemplateService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.templates.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "asset template not found")
		}
		return err
	}
	return nil
}

// ============================================================
// assetService
// ============================================================

type assetService struct {
	db           *sqlx.DB
	assets       repository.AssetRepository
	values       repository.AssetValueRepository
	images       repository.AssetImageRepository
	documents    repository.AssetDocumentRepository
	availability repository.AssetAvailabilityRepository
	fields       repository.TemplateFieldRepository
}

func NewAssetService(
	db *sqlx.DB,
	assets repository.AssetRepository,
	values repository.AssetValueRepository,
	images repository.AssetImageRepository,
	documents repository.AssetDocumentRepository,
	availability repository.AssetAvailabilityRepository,
	fields repository.TemplateFieldRepository,
) AssetService {
	return &assetService{
		db: db, assets: assets, values: values,
		images: images, documents: documents,
		availability: availability, fields: fields,
	}
}

func (s *assetService) Create(ctx context.Context, tenantID, actorID string, req request.CreateAsset) (*response.Asset, error) {
	a := &entity.Asset{
		ID:               uuid.NewString(),
		TenantID:         tenantID,
		CategoryID:       req.CategoryID,
		AssetTemplateID:  req.AssetTemplateID,
		Name:             req.Name,
		SKU:              req.SKU,
		SerialNumber:     req.SerialNumber,
		Description:      req.Description,
		PurchasePrice:    req.PurchasePrice,
		ReplacementValue: req.ReplacementValue,
		PurchaseDate:     req.PurchaseDate,
		Condition:        req.Condition,
		Location:         req.Location,
		CreatedBy:        &actorID,
	}

	if err := transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if err := s.assets.Create(ctx, tx, a); err != nil {
			return err
		}
		for _, v := range req.Values {
			av := &entity.AssetValue{
				ID:              uuid.NewString(),
				TenantID:        tenantID,
				AssetID:         a.ID,
				TemplateFieldID: v.TemplateFieldID,
				Value:           v.Value,
			}
			if err := s.values.Upsert(ctx, tx, av); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, a.ID, tenantID)
}

func (s *assetService) GetByID(ctx context.Context, id, tenantID string) (*response.Asset, error) {
	a, err := s.assets.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "asset not found")
		}
		return nil, err
	}
	return s.hydrate(ctx, a)
}

func (s *assetService) List(ctx context.Context, tenantID string, filter request.ListAssetsFilter) ([]response.Asset, error) {
	perPage := filter.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	page := filter.Page
	if page <= 0 {
		page = 1
	}

	assets, err := s.assets.List(ctx, s.db, tenantID, filter.CategoryID, filter.Status, filter.Search, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}

	out := make([]response.Asset, 0, len(assets))
	for _, a := range assets {
		r, err := s.hydrate(ctx, &a)
		if err != nil {
			return nil, err
		}
		out = append(out, *r)
	}
	return out, nil
}

func (s *assetService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateAsset) (*response.Asset, error) {
	a := &entity.Asset{
		ID:               id,
		TenantID:         tenantID,
		Name:             *req.Name,
		Description:      req.Description,
		Condition:        derefStr(req.Condition),
		Location:         req.Location,
		PurchasePrice:    req.PurchasePrice,
		ReplacementValue: req.ReplacementValue,
		UpdatedBy:        &actorID,
	}
	if err := s.assets.Update(ctx, s.db, a); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "asset not found")
		}
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

func (s *assetService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.assets.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "asset not found")
		}
		return err
	}
	return nil
}

func (s *assetService) AddImage(ctx context.Context, tenantID, assetID, actorID string, img *entity.AssetImage) (*entity.AssetImage, error) {
	img.ID = uuid.NewString()
	img.TenantID = tenantID
	img.AssetID = assetID
	img.CreatedBy = &actorID
	if err := s.images.Create(ctx, s.db, img); err != nil {
		return nil, err
	}
	return img, nil
}

func (s *assetService) DeleteImage(ctx context.Context, tenantID, assetID, imageID, actorID string) error {
	if err := s.images.Delete(ctx, s.db, imageID, assetID, actorID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "image not found")
		}
		return err
	}
	return nil
}

func (s *assetService) AddDocument(ctx context.Context, tenantID, assetID, actorID string, doc *entity.AssetDocument) (*entity.AssetDocument, error) {
	doc.ID = uuid.NewString()
	doc.TenantID = tenantID
	doc.AssetID = assetID
	doc.CreatedBy = &actorID
	if err := s.documents.Create(ctx, s.db, doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *assetService) DeleteDocument(ctx context.Context, tenantID, assetID, docID, actorID string) error {
	if err := s.documents.Delete(ctx, s.db, docID, assetID, actorID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "document not found")
		}
		return err
	}
	return nil
}

func (s *assetService) BlockAvailability(ctx context.Context, tenantID, assetID, actorID string, req request.BlockAvailability) (*entity.AssetAvailability, error) {
	// Guard: asset must exist in this tenant.
	if _, err := s.assets.FindByID(ctx, s.db, assetID, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "asset not found")
		}
		return nil, err
	}

	av := &entity.AssetAvailability{
		ID:               uuid.NewString(),
		TenantID:         tenantID,
		AssetID:          assetID,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		AvailabilityType: req.AvailabilityType,
		Reason:           req.Reason,
		CreatedBy:        &actorID,
	}
	if err := s.availability.Create(ctx, s.db, av); err != nil {
		return nil, err
	}
	return av, nil
}

func (s *assetService) ListAvailability(ctx context.Context, assetID, tenantID string) ([]entity.AssetAvailability, error) {
	return s.availability.ListByAsset(ctx, s.db, assetID, tenantID)
}

// CheckAvailability is the critical cross-module contract. It returns true
// when there are no blocking availability records overlapping [start, end).
// The booking module calls this before reserving any asset.
func (s *assetService) CheckAvailability(ctx context.Context, assetID string, start, end time.Time) (bool, error) {
	count, err := s.availability.CountConflicts(ctx, s.db, assetID, start, end)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// ============================================================
// hydrate — load related data for a single asset
// ============================================================

func (s *assetService) hydrate(ctx context.Context, a *entity.Asset) (*response.Asset, error) {
	vals, _ := s.values.ListByAsset(ctx, s.db, a.ID, a.TenantID)
	imgs, _ := s.images.ListByAsset(ctx, s.db, a.ID, a.TenantID)
	docs, _ := s.documents.ListByAsset(ctx, s.db, a.ID, a.TenantID)

	// Build value responses — include field label/type if template exists.
	var fieldMap map[string]entity.TemplateField
	if a.AssetTemplateID != nil {
		fields, err := s.fields.ListByTemplate(ctx, s.db, *a.AssetTemplateID, a.TenantID)
		if err == nil {
			fieldMap = make(map[string]entity.TemplateField, len(fields))
			for _, f := range fields {
				fieldMap[f.ID] = f
			}
		}
	}

	valueResp := make([]response.AssetValue, 0, len(vals))
	for _, v := range vals {
		av := response.AssetValue{TemplateFieldID: v.TemplateFieldID, Value: v.Value}
		if f, ok := fieldMap[v.TemplateFieldID]; ok {
			av.FieldLabel = f.FieldLabel
			av.FieldType = f.FieldType
		}
		valueResp = append(valueResp, av)
	}

	imgResp := make([]response.AssetImage, 0, len(imgs))
	for _, img := range imgs {
		imgResp = append(imgResp, response.AssetImage{
			ID: img.ID, URL: img.URL, AltText: img.AltText,
			IsPrimary: img.IsPrimary, SortOrder: img.SortOrder,
		})
	}

	docResp := make([]response.AssetDocument, 0, len(docs))
	for _, doc := range docs {
		docResp = append(docResp, response.AssetDocument{
			ID: doc.ID, Title: doc.Title, FileURL: doc.FileURL,
			FileType: doc.FileType, FileSize: doc.FileSize,
		})
	}

	return &response.Asset{
		ID: a.ID, TenantID: a.TenantID, CategoryID: a.CategoryID,
		AssetTemplateID: a.AssetTemplateID, Name: a.Name, SKU: a.SKU,
		SerialNumber: a.SerialNumber, Description: a.Description,
		PurchasePrice: a.PurchasePrice, ReplacementValue: a.ReplacementValue,
		PurchaseDate: a.PurchaseDate, Condition: a.Condition, Location: a.Location,
		Status: a.Status, Values: valueResp, Images: imgResp, Documents: docResp,
		CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
	}, nil
}

// ============================================================
// helpers
// ============================================================

func toTemplateResponse(t *entity.AssetTemplate, fields []entity.TemplateField) *response.AssetTemplate {
	fr := make([]response.TemplateField, 0, len(fields))
	for _, f := range fields {
		fr = append(fr, response.TemplateField{
			ID: f.ID, FieldName: f.FieldName, FieldLabel: f.FieldLabel,
			FieldType: f.FieldType, IsRequired: f.IsRequired,
			DefaultValue: f.DefaultValue, Options: f.Options, SortOrder: f.SortOrder,
		})
	}
	return &response.AssetTemplate{
		ID: t.ID, TenantID: t.TenantID, CategoryID: t.CategoryID,
		Name: t.Name, Description: t.Description, Status: t.Status,
		Fields: fr, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
