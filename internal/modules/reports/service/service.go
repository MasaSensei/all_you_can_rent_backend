package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/reports/dto/request"
	"rentos-backend/internal/modules/reports/dto/response"
	"rentos-backend/internal/modules/reports/entity"
	"rentos-backend/internal/modules/reports/repository"
	pkgresponse "rentos-backend/pkg/response"
)

// ============================================================
// reportService
// ============================================================

type reportService struct {
	db   *sqlx.DB
	repo repository.ReportRepository
}

func NewReportService(db *sqlx.DB, repo repository.ReportRepository) ReportService {
	return &reportService{db: db, repo: repo}
}

func (s *reportService) Generate(ctx context.Context, tenantID, actorID string, req request.GenerateReport) (*response.Report, error) {
	// Serialize parameters to JSON.
	var paramsJSON *string
	if len(req.Parameters) > 0 {
		b, err := json.Marshal(req.Parameters)
		if err != nil {
			return nil, err
		}
		s := string(b)
		paramsJSON = &s
	}

	rep := &entity.Report{
		ID:              uuid.NewString(),
		TenantID:        tenantID,
		Name:            req.Name,
		ReportType:      req.ReportType,
		Parameters:      paramsJSON,
		GeneratedFormat: req.GeneratedFormat,
		CreatedBy:       &actorID,
	}

	if err := s.repo.Create(ctx, s.db, rep); err != nil {
		return nil, err
	}

	// TODO Phase 12: enqueue rep.ID to the report_generator worker queue.
	// The worker will call UpdateStatus once generation completes.

	return s.GetByID(ctx, rep.ID, tenantID)
}

func (s *reportService) GetByID(ctx context.Context, id, tenantID string) (*response.Report, error) {
	rep, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "report not found")
		}
		return nil, err
	}
	return toReportResponse(rep), nil
}

func (s *reportService) List(ctx context.Context, tenantID string, page, perPage int) ([]response.Report, error) {
	perPage, page = normPage(perPage, page)
	reports, err := s.repo.List(ctx, s.db, tenantID, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.Report, 0, len(reports))
	for _, r := range reports {
		out = append(out, *toReportResponse(&r))
	}
	return out, nil
}

func (s *reportService) UpdateStatus(ctx context.Context, id, tenantID, status string, fileURL *string) error {
	return s.repo.UpdateStatus(ctx, s.db, id, tenantID, status, fileURL)
}

// ============================================================
// analyticsService
// ============================================================

type analyticsService struct {
	db   *sqlx.DB
	repo repository.AnalyticsRepository
}

func NewAnalyticsService(db *sqlx.DB, repo repository.AnalyticsRepository) AnalyticsService {
	return &analyticsService{db: db, repo: repo}
}

func (s *analyticsService) Track(ctx context.Context, tenantID string, userID *string, req request.TrackEvent) error {
	occurredAt := time.Now()
	if req.OccurredAt != nil {
		occurredAt = *req.OccurredAt
	}

	var dataJSON *string
	if len(req.EventData) > 0 {
		b, err := json.Marshal(req.EventData)
		if err != nil {
			return err
		}
		s := string(b)
		dataJSON = &s
	}

	e := &entity.AnalyticsEvent{
		ID:            uuid.NewString(),
		TenantID:      tenantID,
		UserID:        userID,
		CustomerID:    req.CustomerID,
		EventName:     req.EventName,
		EventCategory: req.EventCategory,
		EventData:     dataJSON,
		Source:        req.Source,
		OccurredAt:    occurredAt,
	}
	return s.repo.CreateEvent(ctx, s.db, e)
}

func (s *analyticsService) Dashboard(ctx context.Context, tenantID string, filter request.DashboardFilter) (*response.Dashboard, error) {
	from, to := filter.From, filter.To

	totalRevenue, err := s.repo.SumTotalRevenue(ctx, s.db, tenantID, from, to)
	if err != nil {
		return nil, err
	}

	totalBookings, err := s.repo.CountTotalBookings(ctx, s.db, tenantID, from, to)
	if err != nil {
		return nil, err
	}

	activeCustomers, err := s.repo.CountActiveCustomers(ctx, s.db, tenantID, from, to)
	if err != nil {
		return nil, err
	}

	revenueRows, err := s.repo.AggregateRevenue(ctx, s.db, tenantID, from, to, filter.GroupBy)
	if err != nil {
		return nil, err
	}
	revenueSeries := make([]response.RevenueDataPoint, 0, len(revenueRows))
	for _, r := range revenueRows {
		revenueSeries = append(revenueSeries, response.RevenueDataPoint{
			Period:  r.Period,
			Revenue: r.Revenue,
			Count:   r.Count,
		})
	}

	utilRows, err := s.repo.AggregateAssetUtilization(ctx, s.db, tenantID, from, to, 10)
	if err != nil {
		return nil, err
	}
	topAssets := make([]response.AssetUtilizationItem, 0, len(utilRows))
	for _, u := range utilRows {
		topAssets = append(topAssets, response.AssetUtilizationItem{
			AssetID:     u.AssetID,
			AssetName:   u.AssetName,
			BookedDays:  u.BookedDays,
			TotalDays:   u.TotalDays,
			Utilization: u.Utilization,
		})
	}

	return &response.Dashboard{
		TotalRevenue:    totalRevenue,
		TotalBookings:   totalBookings,
		ActiveCustomers: activeCustomers,
		Revenue:         revenueSeries,
		TopAssets:       topAssets,
	}, nil
}

// ============================================================
// helpers
// ============================================================

func toReportResponse(r *entity.Report) *response.Report {
	return &response.Report{
		ID:              r.ID,
		TenantID:        r.TenantID,
		Name:            r.Name,
		ReportType:      r.ReportType,
		GeneratedFormat: r.GeneratedFormat,
		FileURL:         r.FileURL,
		Status:          r.Status,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
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
