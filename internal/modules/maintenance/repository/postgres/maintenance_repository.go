package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"rentos-backend/internal/modules/maintenance/entity"
	"rentos-backend/internal/modules/maintenance/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// maintenanceRepository
// ============================================================

type maintenanceRepository struct {
	qCreate       string
	qFindByID     string
	qList         string
	qListDue      string
	qUpdateStatus string
	qDelete       string
}

func NewMaintenanceRepository(qCreate, qFindByID, qList, qListDue, qUpdateStatus, qDelete string) repository.MaintenanceRepository {
	return &maintenanceRepository{
		qCreate: qCreate, qFindByID: qFindByID, qList: qList,
		qListDue: qListDue, qUpdateStatus: qUpdateStatus, qDelete: qDelete,
	}
}

func (r *maintenanceRepository) Create(ctx context.Context, q database.Querier, rec *entity.MaintenanceRecord) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		rec.ID, rec.TenantID, rec.AssetID, rec.MaintenanceType, rec.Description,
		rec.Cost, rec.ScheduledDate, rec.PerformedBy, rec.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("maintenanceRepository.Create: %w", err)
	}
	return nil
}

func (r *maintenanceRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.MaintenanceRecord, error) {
	var rec entity.MaintenanceRecord
	if err := q.GetContext(ctx, &rec, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("maintenanceRepository.FindByID: %w", err)
	}
	return &rec, nil
}

func (r *maintenanceRepository) List(ctx context.Context, q database.Querier, tenantID string, assetID, status *string, limit, offset int) ([]entity.MaintenanceRecord, error) {
	var out []entity.MaintenanceRecord
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, assetID, status, limit, offset); err != nil {
		return nil, fmt.Errorf("maintenanceRepository.List: %w", err)
	}
	return out, nil
}

func (r *maintenanceRepository) ListDue(ctx context.Context, q database.Querier, dueBy time.Time) ([]entity.MaintenanceRecord, error) {
	var out []entity.MaintenanceRecord
	if err := q.SelectContext(ctx, &out, r.qListDue, dueBy); err != nil {
		return nil, fmt.Errorf("maintenanceRepository.ListDue: %w", err)
	}
	return out, nil
}

func (r *maintenanceRepository) UpdateStatus(ctx context.Context, q database.Querier, rec *entity.MaintenanceRecord) error {
	res, err := q.ExecContext(ctx, r.qUpdateStatus,
		rec.ID, rec.TenantID, rec.MaintenanceStatus,
		rec.CompletedDate, rec.Cost, rec.PerformedBy,
	)
	if err != nil {
		return fmt.Errorf("maintenanceRepository.UpdateStatus: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *maintenanceRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("maintenanceRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// inspectionRepository
// ============================================================

type inspectionRepository struct {
	qCreate   string
	qFindByID string
	qList     string
}

func NewInspectionRepository(qCreate, qFindByID, qList string) repository.InspectionRepository {
	return &inspectionRepository{qCreate: qCreate, qFindByID: qFindByID, qList: qList}
}

func (r *inspectionRepository) Create(ctx context.Context, q database.Querier, i *entity.Inspection) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		i.ID, i.TenantID, i.AssetID, i.BookingItemID, i.InspectionType,
		i.InspectorName, i.Findings, i.Result, i.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("inspectionRepository.Create: %w", err)
	}
	return nil
}

func (r *inspectionRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Inspection, error) {
	var i entity.Inspection
	if err := q.GetContext(ctx, &i, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("inspectionRepository.FindByID: %w", err)
	}
	return &i, nil
}

func (r *inspectionRepository) List(ctx context.Context, q database.Querier, tenantID string, assetID, inspectionType, result *string, limit, offset int) ([]entity.Inspection, error) {
	var out []entity.Inspection
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, assetID, inspectionType, result, limit, offset); err != nil {
		return nil, fmt.Errorf("inspectionRepository.List: %w", err)
	}
	return out, nil
}

// ============================================================
// damageReportRepository
// ============================================================

type damageReportRepository struct {
	qCreate       string
	qFindByID     string
	qList         string
	qUpdateStatus string
}

func NewDamageReportRepository(qCreate, qFindByID, qList, qUpdateStatus string) repository.DamageReportRepository {
	return &damageReportRepository{
		qCreate: qCreate, qFindByID: qFindByID,
		qList: qList, qUpdateStatus: qUpdateStatus,
	}
}

func (r *damageReportRepository) Create(ctx context.Context, q database.Querier, rep *entity.DamageReport) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		rep.ID, rep.TenantID, rep.AssetID, rep.BookingID, rep.InspectionID,
		rep.Description, rep.Severity, rep.RepairCost, rep.ChargedAmount, rep.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("damageReportRepository.Create: %w", err)
	}
	return nil
}

func (r *damageReportRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.DamageReport, error) {
	var rep entity.DamageReport
	if err := q.GetContext(ctx, &rep, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("damageReportRepository.FindByID: %w", err)
	}
	return &rep, nil
}

func (r *damageReportRepository) List(ctx context.Context, q database.Querier, tenantID string, assetID, status, severity *string, limit, offset int) ([]entity.DamageReport, error) {
	var out []entity.DamageReport
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, assetID, status, severity, limit, offset); err != nil {
		return nil, fmt.Errorf("damageReportRepository.List: %w", err)
	}
	return out, nil
}

func (r *damageReportRepository) UpdateStatus(ctx context.Context, q database.Querier, rep *entity.DamageReport) error {
	res, err := q.ExecContext(ctx, r.qUpdateStatus,
		rep.ID, rep.TenantID, rep.ReportStatus, rep.RepairCost, rep.ChargedAmount,
	)
	if err != nil {
		return fmt.Errorf("damageReportRepository.UpdateStatus: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}
