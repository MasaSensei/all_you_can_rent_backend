package service

import (
	"context"

	"rentos-backend/internal/modules/reports/dto/request"
	"rentos-backend/internal/modules/reports/dto/response"
)

// ReportService manages report generation jobs.
type ReportService interface {
	Generate(ctx context.Context, tenantID, actorID string, req request.GenerateReport) (*response.Report, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Report, error)
	List(ctx context.Context, tenantID string, page, perPage int) ([]response.Report, error)
	// UpdateStatus is called by the background report generator worker.
	UpdateStatus(ctx context.Context, id, tenantID, status string, fileURL *string) error
}

// AnalyticsService manages event ingestion and dashboard aggregation.
type AnalyticsService interface {
	Track(ctx context.Context, tenantID string, userID *string, req request.TrackEvent) error
	Dashboard(ctx context.Context, tenantID string, filter request.DashboardFilter) (*response.Dashboard, error)
}
