package service

import (
	"context"
	"time"

	"rentos-backend/internal/modules/maintenance/dto/request"
	"rentos-backend/internal/modules/maintenance/dto/response"
)

// MaintenanceService manages maintenance records lifecycle.
type MaintenanceService interface {
	Schedule(ctx context.Context, tenantID, actorID string, req request.ScheduleMaintenance) (*response.MaintenanceRecord, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.MaintenanceRecord, error)
	List(ctx context.Context, tenantID string, filter request.ListMaintenanceFilter) ([]response.MaintenanceRecord, error)
	UpdateStatus(ctx context.Context, id, tenantID, actorID string, req request.UpdateMaintenanceStatus) (*response.MaintenanceRecord, error)
	Delete(ctx context.Context, id, tenantID string) error
	// ListDue is called by the maintenance scheduler background worker.
	ListDue(ctx context.Context, dueBy time.Time) ([]response.MaintenanceRecord, error)
}

// InspectionService manages asset inspections.
type InspectionService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateInspection) (*response.Inspection, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Inspection, error)
	List(ctx context.Context, tenantID string, filter request.ListInspectionsFilter) ([]response.Inspection, error)
}

// DamageReportService manages damage reports.
type DamageReportService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateDamageReport) (*response.DamageReport, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.DamageReport, error)
	List(ctx context.Context, tenantID string, filter request.ListDamageReportsFilter) ([]response.DamageReport, error)
	UpdateStatus(ctx context.Context, id, tenantID, actorID string, req request.UpdateDamageReportStatus) (*response.DamageReport, error)
}
