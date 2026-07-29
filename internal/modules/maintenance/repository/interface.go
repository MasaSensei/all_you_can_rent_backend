package repository

import (
	"context"
	"errors"
	"time"

	"rentos-backend/internal/modules/maintenance/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

// MaintenanceRepository manages the maintenance_records table.
type MaintenanceRepository interface {
	Create(ctx context.Context, q database.Querier, r *entity.MaintenanceRecord) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.MaintenanceRecord, error)
	List(ctx context.Context, q database.Querier, tenantID string, assetID, status *string, limit, offset int) ([]entity.MaintenanceRecord, error)
	// ListDue returns all scheduled records whose scheduled_date <= dueBy.
	// Used by the maintenance scheduler background worker.
	ListDue(ctx context.Context, q database.Querier, dueBy time.Time) ([]entity.MaintenanceRecord, error)
	UpdateStatus(ctx context.Context, q database.Querier, r *entity.MaintenanceRecord) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// InspectionRepository manages the inspections table.
type InspectionRepository interface {
	Create(ctx context.Context, q database.Querier, i *entity.Inspection) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Inspection, error)
	List(ctx context.Context, q database.Querier, tenantID string, assetID, inspectionType, result *string, limit, offset int) ([]entity.Inspection, error)
}

// DamageReportRepository manages the damage_reports table.
type DamageReportRepository interface {
	Create(ctx context.Context, q database.Querier, r *entity.DamageReport) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.DamageReport, error)
	List(ctx context.Context, q database.Querier, tenantID string, assetID, status, severity *string, limit, offset int) ([]entity.DamageReport, error)
	UpdateStatus(ctx context.Context, q database.Querier, r *entity.DamageReport) error
}
