package handler

import (
	"github.com/gofiber/fiber/v2"

	rreq "rentos-backend/internal/modules/reports/dto/request"
	"rentos-backend/internal/modules/reports/service"
	apiresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

const (
	ctxKeyTenantID = "tenant_id"
	ctxKeyUserID   = "user_id"
)

// Handler groups all reports and analytics HTTP handlers.
type Handler struct {
	reports   service.ReportService
	analytics service.AnalyticsService
	validate  *validator.Validate
}

func New(reports service.ReportService, analytics service.AnalyticsService, v *validator.Validate) *Handler {
	return &Handler{reports: reports, analytics: analytics, validate: v}
}

func tenantID(c *fiber.Ctx) string {
	if id, ok := c.Locals(ctxKeyTenantID).(string); ok {
		return id
	}
	return c.Get("X-Tenant-ID")
}

func userID(c *fiber.Ctx) string {
	id, _ := c.Locals(ctxKeyUserID).(string)
	return id
}

// ---- Reports ----

func (h *Handler) GenerateReport(c *fiber.Ctx) error {
	var req rreq.GenerateReport
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	rep, err := h.reports.Generate(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, rep)
}

func (h *Handler) GetReport(c *fiber.Ctx) error {
	rep, err := h.reports.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, rep)
}

func (h *Handler) ListReports(c *fiber.Ctx) error {
	reports, err := h.reports.List(c.Context(), tenantID(c),
		c.QueryInt("page", 1), c.QueryInt("per_page", 20),
	)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, reports)
}

// ---- Analytics ----

func (h *Handler) TrackEvent(c *fiber.Ctx) error {
	var req rreq.TrackEvent
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	uid := userID(c)
	var uIDPtr *string
	if uid != "" {
		uIDPtr = &uid
	}

	if err := h.analytics.Track(c.Context(), tenantID(c), uIDPtr, req); err != nil {
		return apiresponse.FromError(c, err)
	}
	// 202 Accepted — event is recorded but no body needed.
	return c.SendStatus(fiber.StatusAccepted)
}

func (h *Handler) Dashboard(c *fiber.Ctx) error {
	var filter rreq.DashboardFilter
	if err := c.BodyParser(&filter); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(filter); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	dashboard, err := h.analytics.Dashboard(c.Context(), tenantID(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, dashboard)
}
