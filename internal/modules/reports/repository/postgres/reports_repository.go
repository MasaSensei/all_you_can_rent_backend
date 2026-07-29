package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"rentos-backend/internal/modules/reports/entity"
	"rentos-backend/internal/modules/reports/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// reportRepository
// ============================================================

type reportRepository struct {
	qCreate       string
	qFindByID     string
	qList         string
	qUpdateStatus string
}

func NewReportRepository(qCreate, qFindByID, qList, qUpdateStatus string) repository.ReportRepository {
	return &reportRepository{
		qCreate: qCreate, qFindByID: qFindByID,
		qList: qList, qUpdateStatus: qUpdateStatus,
	}
}

func (r *reportRepository) Create(ctx context.Context, q database.Querier, rep *entity.Report) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		rep.ID, rep.TenantID, rep.Name, rep.ReportType,
		rep.Parameters, rep.GeneratedFormat, rep.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("reportRepository.Create: %w", err)
	}
	return nil
}

func (r *reportRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Report, error) {
	var rep entity.Report
	if err := q.GetContext(ctx, &rep, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("reportRepository.FindByID: %w", err)
	}
	return &rep, nil
}

func (r *reportRepository) List(ctx context.Context, q database.Querier, tenantID string, limit, offset int) ([]entity.Report, error) {
	var out []entity.Report
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, limit, offset); err != nil {
		return nil, fmt.Errorf("reportRepository.List: %w", err)
	}
	return out, nil
}

func (r *reportRepository) UpdateStatus(ctx context.Context, q database.Querier, id, tenantID, status string, fileURL *string) error {
	_, err := q.ExecContext(ctx, r.qUpdateStatus, id, tenantID, status, fileURL)
	if err != nil {
		return fmt.Errorf("reportRepository.UpdateStatus: %w", err)
	}
	return nil
}

// ============================================================
// analyticsRepository
// ============================================================

type analyticsRepository struct {
	qCreateEvent           string
	qAggregateRevenue      string
	qAggregateUtilization  string
	qCountCustomers        string
	qCountBookings         string
	qSumRevenue            string
}

func NewAnalyticsRepository(
	qCreateEvent, qAggregateRevenue, qAggregateUtilization,
	qCountCustomers, qCountBookings, qSumRevenue string,
) repository.AnalyticsRepository {
	return &analyticsRepository{
		qCreateEvent:          qCreateEvent,
		qAggregateRevenue:     qAggregateRevenue,
		qAggregateUtilization: qAggregateUtilization,
		qCountCustomers:       qCountCustomers,
		qCountBookings:        qCountBookings,
		qSumRevenue:           qSumRevenue,
	}
}

func (r *analyticsRepository) CreateEvent(ctx context.Context, q database.Querier, e *entity.AnalyticsEvent) error {
	_, err := q.ExecContext(ctx, r.qCreateEvent,
		e.ID, e.TenantID, e.UserID, e.CustomerID,
		e.EventName, e.EventCategory, e.EventData, e.Source,
		e.OccurredAt,
	)
	if err != nil {
		return fmt.Errorf("analyticsRepository.CreateEvent: %w", err)
	}
	return nil
}

func (r *analyticsRepository) AggregateRevenue(ctx context.Context, q database.Querier, tenantID string, from, to time.Time, groupBy string) ([]repository.RevenueRow, error) {
	var out []repository.RevenueRow
	if err := q.SelectContext(ctx, &out, r.qAggregateRevenue, tenantID, from, to, groupBy); err != nil {
		return nil, fmt.Errorf("analyticsRepository.AggregateRevenue: %w", err)
	}
	return out, nil
}

func (r *analyticsRepository) AggregateAssetUtilization(ctx context.Context, q database.Querier, tenantID string, from, to time.Time, topN int) ([]repository.UtilizationRow, error) {
	var out []repository.UtilizationRow
	if err := q.SelectContext(ctx, &out, r.qAggregateUtilization, tenantID, from, to, topN); err != nil {
		return nil, fmt.Errorf("analyticsRepository.AggregateAssetUtilization: %w", err)
	}
	return out, nil
}

func (r *analyticsRepository) CountActiveCustomers(ctx context.Context, q database.Querier, tenantID string, from, to time.Time) (int, error) {
	var count int
	if err := q.GetContext(ctx, &count, r.qCountCustomers, tenantID, from, to); err != nil {
		return 0, fmt.Errorf("analyticsRepository.CountActiveCustomers: %w", err)
	}
	return count, nil
}

func (r *analyticsRepository) CountTotalBookings(ctx context.Context, q database.Querier, tenantID string, from, to time.Time) (int, error) {
	var count int
	if err := q.GetContext(ctx, &count, r.qCountBookings, tenantID, from, to); err != nil {
		return 0, fmt.Errorf("analyticsRepository.CountTotalBookings: %w", err)
	}
	return count, nil
}

func (r *analyticsRepository) SumTotalRevenue(ctx context.Context, q database.Querier, tenantID string, from, to time.Time) (float64, error) {
	var total float64
	if err := q.GetContext(ctx, &total, r.qSumRevenue, tenantID, from, to); err != nil {
		return 0, fmt.Errorf("analyticsRepository.SumTotalRevenue: %w", err)
	}
	return total, nil
}
