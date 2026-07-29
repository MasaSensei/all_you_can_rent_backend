package maintenance

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/maintenance/handler"
	"rentos-backend/internal/modules/maintenance/repository/postgres"
	"rentos-backend/internal/modules/maintenance/routes"
	"rentos-backend/internal/modules/maintenance/service"
)

// Module holds the maintenance module's wired handler and services.
type Module struct {
	handler        *handler.Handler
	maintenanceSvc service.MaintenanceService
}

// New builds the maintenance module: repositories → services → handler.
func New(c *bootstrap.Container) *Module {
	maintenanceRepo := postgres.NewMaintenanceRepository(
		query("create_maintenance_record.sql"),
		query("find_maintenance_record_by_id.sql"),
		query("list_maintenance_records.sql"),
		query("list_due_maintenance.sql"),
		query("update_maintenance_status.sql"),
		query("delete_maintenance_record.sql"),
	)
	inspectionRepo := postgres.NewInspectionRepository(
		query("create_inspection.sql"),
		query("find_inspection_by_id.sql"),
		query("list_inspections.sql"),
	)
	damageReportRepo := postgres.NewDamageReportRepository(
		query("create_damage_report.sql"),
		query("find_damage_report_by_id.sql"),
		query("list_damage_reports.sql"),
		query("update_damage_report_status.sql"),
	)

	maintenanceSvc := service.NewMaintenanceService(c.DB, maintenanceRepo)
	inspectionSvc := service.NewInspectionService(c.DB, inspectionRepo)
	damageReportSvc := service.NewDamageReportService(c.DB, damageReportRepo)

	h := handler.New(maintenanceSvc, inspectionSvc, damageReportSvc, c.Validator)
	return &Module{handler: h, maintenanceSvc: maintenanceSvc}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}

// MaintenanceService exposes the service so the background scheduler
// worker can call ListDue without importing internal packages.
func (m *Module) MaintenanceService() service.MaintenanceService {
	return m.maintenanceSvc
}
