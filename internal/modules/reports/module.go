package reports

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/reports/handler"
	"rentos-backend/internal/modules/reports/repository/postgres"
	"rentos-backend/internal/modules/reports/routes"
	"rentos-backend/internal/modules/reports/service"
)

// Module holds the reports module's wired handler and services.
type Module struct {
	handler    *handler.Handler
	reportSvc  service.ReportService
}

// New builds the reports module: repositories → services → handler.
func New(c *bootstrap.Container) *Module {
	reportRepo := postgres.NewReportRepository(
		query("create_report.sql"),
		query("find_report_by_id.sql"),
		query("list_reports.sql"),
		query("update_report_status.sql"),
	)
	analyticsRepo := postgres.NewAnalyticsRepository(
		query("create_analytics_event.sql"),
		query("aggregate_revenue_by_period.sql"),
		query("aggregate_asset_utilization.sql"),
		query("count_active_customers.sql"),
		query("count_total_bookings.sql"),
		query("sum_total_revenue.sql"),
	)

	reportSvc    := service.NewReportService(c.DB, reportRepo)
	analyticsSvc := service.NewAnalyticsService(c.DB, analyticsRepo)

	h := handler.New(reportSvc, analyticsSvc, c.Validator)
	return &Module{handler: h, reportSvc: reportSvc}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}

// ReportService exposes the service so the background report generator
// worker can call UpdateStatus without importing internal packages.
func (m *Module) ReportService() service.ReportService {
	return m.reportSvc
}
