package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/maintenance/handler"
)

// Register mounts all maintenance routes under /api/v1.
func Register(router fiber.Router, h *handler.Handler) {
	maintenance := router.Group("/maintenance-records")
	maintenance.Post("/", h.ScheduleMaintenance)
	maintenance.Get("/", h.ListMaintenanceRecords)
	maintenance.Get("/:id", h.GetMaintenanceRecord)
	maintenance.Patch("/:id/status", h.UpdateMaintenanceStatus)
	maintenance.Delete("/:id", h.DeleteMaintenanceRecord)

	inspections := router.Group("/inspections")
	inspections.Post("/", h.CreateInspection)
	inspections.Get("/", h.ListInspections)
	inspections.Get("/:id", h.GetInspection)

	damageReports := router.Group("/damage-reports")
	damageReports.Post("/", h.CreateDamageReport)
	damageReports.Get("/", h.ListDamageReports)
	damageReports.Get("/:id", h.GetDamageReport)
	damageReports.Patch("/:id/status", h.UpdateDamageReportStatus)
}
