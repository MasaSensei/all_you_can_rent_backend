package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/reports/handler"
)

// Register mounts all report and analytics routes under /api/v1.
func Register(router fiber.Router, h *handler.Handler) {
	reports := router.Group("/reports")
	reports.Post("/", h.GenerateReport)
	reports.Get("/", h.ListReports)
	reports.Get("/:id", h.GetReport)

	analytics := router.Group("/analytics")
	// POST /analytics/events — fire-and-forget event ingestion
	analytics.Post("/events", h.TrackEvent)
	// POST /analytics/dashboard — aggregated KPI dashboard
	analytics.Post("/dashboard", h.Dashboard)
}
