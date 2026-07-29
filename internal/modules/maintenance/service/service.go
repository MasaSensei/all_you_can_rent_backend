package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/maintenance/dto/request"
	"rentos-backend/internal/modules/maintenance/dto/response"
	"rentos-backend/internal/modules/maintenance/entity"
	"rentos-backend/internal/modules/maintenance/repository"
	pkgresponse "rentos-backend/pkg/response"
)

// ============================================================
// maintenanceService
// ============================================================

type maintenanceService struct {
	db   *sqlx.DB
	repo repository.MaintenanceRepository
}

func NewMaintenanceService(db *sqlx.DB, repo repository.MaintenanceRepository) MaintenanceService {
	return &maintenanceService{db: db, repo: repo}
}

func (s *maintenanceService) Schedule(ctx context.Context, tenantID, actorID string, req request.ScheduleMaintenance) (*response.MaintenanceRecord, error) {
	rec := &entity.MaintenanceRecord{
		ID:              uuid.NewString(),
		TenantID:        tenantID,
		AssetID:         req.AssetID,
		MaintenanceType: req.MaintenanceType,
		Description:     req.Description,
		Cost:            req.Cost,
		ScheduledDate:   req.ScheduledDate,
		PerformedBy:     req.PerformedBy,
		CreatedBy:       &actorID,
	}
	if err := s.repo.Create(ctx, s.db, rec); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, rec.ID, tenantID)
}

func (s *maintenanceService) GetByID(ctx context.Context, id, tenantID string) (*response.MaintenanceRecord, error) {
	rec, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "maintenance record not found")
		}
		return nil, err
	}
	return toMaintenanceResponse(rec), nil
}

func (s *maintenanceService) List(ctx context.Context, tenantID string, filter request.ListMaintenanceFilter) ([]response.MaintenanceRecord, error) {
	perPage, page := normPage(filter.PerPage, filter.Page)
	records, err := s.repo.List(ctx, s.db, tenantID, filter.AssetID, filter.MaintenanceStatus, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.MaintenanceRecord, 0, len(records))
	for _, r := range records {
		out = append(out, *toMaintenanceResponse(&r))
	}
	return out, nil
}

func (s *maintenanceService) UpdateStatus(ctx context.Context, id, tenantID, actorID string, req request.UpdateMaintenanceStatus) (*response.MaintenanceRecord, error) {
	rec := &entity.MaintenanceRecord{
		ID:                id,
		TenantID:          tenantID,
		MaintenanceStatus: req.MaintenanceStatus,
		CompletedDate:     req.CompletedDate,
		Cost:              derefFloat(req.Cost),
		PerformedBy:       req.PerformedBy,
		UpdatedBy:         &actorID,
	}
	if err := s.repo.UpdateStatus(ctx, s.db, rec); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "maintenance record not found")
		}
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

func (s *maintenanceService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.repo.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "maintenance record not found")
		}
		return err
	}
	return nil
}

func (s *maintenanceService) ListDue(ctx context.Context, dueBy time.Time) ([]response.MaintenanceRecord, error) {
	records, err := s.repo.ListDue(ctx, s.db, dueBy)
	if err != nil {
		return nil, err
	}
	out := make([]response.MaintenanceRecord, 0, len(records))
	for _, r := range records {
		out = append(out, *toMaintenanceResponse(&r))
	}
	return out, nil
}

// ============================================================
// inspectionService
// ============================================================

type inspectionService struct {
	db   *sqlx.DB
	repo repository.InspectionRepository
}

func NewInspectionService(db *sqlx.DB, repo repository.InspectionRepository) InspectionService {
	return &inspectionService{db: db, repo: repo}
}

func (s *inspectionService) Create(ctx context.Context, tenantID, actorID string, req request.CreateInspection) (*response.Inspection, error) {
	i := &entity.Inspection{
		ID:             uuid.NewString(),
		TenantID:       tenantID,
		AssetID:        req.AssetID,
		BookingItemID:  req.BookingItemID,
		InspectionType: req.InspectionType,
		InspectedAt:    time.Now(),
		InspectorName:  req.InspectorName,
		Findings:       req.Findings,
		Result:         req.Result,
		CreatedBy:      &actorID,
	}
	if err := s.repo.Create(ctx, s.db, i); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, i.ID, tenantID)
}

func (s *inspectionService) GetByID(ctx context.Context, id, tenantID string) (*response.Inspection, error) {
	i, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "inspection not found")
		}
		return nil, err
	}
	return toInspectionResponse(i), nil
}

func (s *inspectionService) List(ctx context.Context, tenantID string, filter request.ListInspectionsFilter) ([]response.Inspection, error) {
	perPage, page := normPage(filter.PerPage, filter.Page)
	items, err := s.repo.List(ctx, s.db, tenantID, filter.AssetID, filter.InspectionType, filter.Result, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.Inspection, 0, len(items))
	for _, i := range items {
		out = append(out, *toInspectionResponse(&i))
	}
	return out, nil
}

// ============================================================
// damageReportService
// ============================================================

type damageReportService struct {
	db   *sqlx.DB
	repo repository.DamageReportRepository
}

func NewDamageReportService(db *sqlx.DB, repo repository.DamageReportRepository) DamageReportService {
	return &damageReportService{db: db, repo: repo}
}

func (s *damageReportService) Create(ctx context.Context, tenantID, actorID string, req request.CreateDamageReport) (*response.DamageReport, error) {
	rep := &entity.DamageReport{
		ID:            uuid.NewString(),
		TenantID:      tenantID,
		AssetID:       req.AssetID,
		BookingID:     req.BookingID,
		InspectionID:  req.InspectionID,
		Description:   req.Description,
		Severity:      req.Severity,
		RepairCost:    req.RepairCost,
		ChargedAmount: req.ChargedAmount,
		CreatedBy:     &actorID,
	}
	if err := s.repo.Create(ctx, s.db, rep); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, rep.ID, tenantID)
}

func (s *damageReportService) GetByID(ctx context.Context, id, tenantID string) (*response.DamageReport, error) {
	rep, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "damage report not found")
		}
		return nil, err
	}
	return toDamageReportResponse(rep), nil
}

func (s *damageReportService) List(ctx context.Context, tenantID string, filter request.ListDamageReportsFilter) ([]response.DamageReport, error) {
	perPage, page := normPage(filter.PerPage, filter.Page)
	items, err := s.repo.List(ctx, s.db, tenantID, filter.AssetID, filter.ReportStatus, filter.Severity, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.DamageReport, 0, len(items))
	for _, r := range items {
		out = append(out, *toDamageReportResponse(&r))
	}
	return out, nil
}

func (s *damageReportService) UpdateStatus(ctx context.Context, id, tenantID, actorID string, req request.UpdateDamageReportStatus) (*response.DamageReport, error) {
	rep := &entity.DamageReport{
		ID:            id,
		TenantID:      tenantID,
		ReportStatus:  req.ReportStatus,
		RepairCost:    derefFloat(req.RepairCost),
		ChargedAmount: derefFloat(req.ChargedAmount),
		UpdatedBy:     &actorID,
	}
	if err := s.repo.UpdateStatus(ctx, s.db, rep); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "damage report not found")
		}
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

// ============================================================
// mapping helpers
// ============================================================

func toMaintenanceResponse(r *entity.MaintenanceRecord) *response.MaintenanceRecord {
	return &response.MaintenanceRecord{
		ID: r.ID, TenantID: r.TenantID, AssetID: r.AssetID,
		MaintenanceType: r.MaintenanceType, Description: r.Description,
		Cost: r.Cost, ScheduledDate: r.ScheduledDate, CompletedDate: r.CompletedDate,
		PerformedBy: r.PerformedBy, MaintenanceStatus: r.MaintenanceStatus,
		Status: r.Status, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func toInspectionResponse(i *entity.Inspection) *response.Inspection {
	return &response.Inspection{
		ID: i.ID, TenantID: i.TenantID, AssetID: i.AssetID,
		BookingItemID: i.BookingItemID, InspectionType: i.InspectionType,
		InspectedAt: i.InspectedAt, InspectorName: i.InspectorName,
		Findings: i.Findings, Result: i.Result, Status: i.Status,
		CreatedAt: i.CreatedAt, UpdatedAt: i.UpdatedAt,
	}
}

func toDamageReportResponse(r *entity.DamageReport) *response.DamageReport {
	return &response.DamageReport{
		ID: r.ID, TenantID: r.TenantID, AssetID: r.AssetID,
		BookingID: r.BookingID, InspectionID: r.InspectionID,
		Description: r.Description, Severity: r.Severity,
		RepairCost: r.RepairCost, ChargedAmount: r.ChargedAmount,
		ReportStatus: r.ReportStatus, Status: r.Status,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
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

func derefFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
