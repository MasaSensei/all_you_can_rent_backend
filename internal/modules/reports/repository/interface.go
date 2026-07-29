package repository

import (
	"context"
	"errors"
	"time"

	"rentos-backend/internal/modules/reports/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

// ReportRepository manages the reports table.
type ReportRepository interface {
	Create(ctx context.Context, q database.Querier, r *entity.Report) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Report, error)
	List(ctx context.Context, q database.Querier, tenantID string, limit, offset int) ([]entity.Report, error)
	// UpdateStatus is called by the report generator worker after completing
	// or failing generation.
	UpdateStatus(ctx context.Context, q database.Querier, id, tenantID, status string, fileURL *string) error
}

// AnalyticsRepository manages analytics_events and dashboard aggregates.
type AnalyticsRepository interface {
	CreateEvent(ctx context.Context, q database.Querier, e *entity.AnalyticsEvent) error

	// AggregateRevenue returns revenue bucketed by day/week/month.
	AggregateRevenue(ctx context.Context, q database.Querier, tenantID string, from, to time.Time, groupBy string) ([]RevenueRow, error)

	// AggregateAssetUtilization returns the top N most-booked assets.
	AggregateAssetUtilization(ctx context.Context, q database.Querier, tenantID string, from, to time.Time, topN int) ([]UtilizationRow, error)

	// CountActiveCustomers returns distinct customers with bookings in range.
	CountActiveCustomers(ctx context.Context, q database.Querier, tenantID string, from, to time.Time) (int, error)

	// CountTotalBookings returns total bookings in range.
	CountTotalBookings(ctx context.Context, q database.Querier, tenantID string, from, to time.Time) (int, error)

	// SumTotalRevenue returns succeeded payment total in range.
	SumTotalRevenue(ctx context.Context, q database.Querier, tenantID string, from, to time.Time) (float64, error)
}

// RevenueRow is the raw DB row from aggregate_revenue_by_period.sql.
type RevenueRow struct {
	Period  string  `db:"period"`
	Revenue float64 `db:"revenue"`
	Count   int     `db:"count"`
}

// UtilizationRow is the raw DB row from aggregate_asset_utilization.sql.
type UtilizationRow struct {
	AssetID     string  `db:"asset_id"`
	AssetName   string  `db:"asset_name"`
	BookedDays  int     `db:"booked_days"`
	TotalDays   int     `db:"total_days"`
	Utilization float64 `db:"utilization_pct"`
}
