package handler

import (
	"github.com/gofiber/fiber/v2"

	mreq "rentos-backend/internal/modules/maintenance/dto/request"
	"rentos-backend/internal/modules/maintenance/service"
	apiresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

const (
	ctxKeyTenantID = "tenant_id"
	ctxKeyUserID   = "user_id"
)

// Handler groups all maintenance HTTP handlers.
type Handler struct {
	maintenance  service.MaintenanceService
	inspections  service.InspectionService
	damageReports service.DamageReportService
	validate     *validator.Validate
}

func New(
	maintenance service.MaintenanceService,
	inspections service.InspectionService,
	damageReports service.DamageReportService,
	v *validator.Validate,
) *Handler {
	return &Handler{
		maintenance: maintenance, inspections: inspections,
		damageReports: damageReports, validate: v,
	}
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

// ---- Maintenance Records ----

func (h *Handler) ScheduleMaintenance(c *fiber.Ctx) error {
	var req mreq.ScheduleMaintenance
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	rec, err := h.maintenance.Schedule(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, rec)
}

func (h *Handler) GetMaintenanceRecord(c *fiber.Ctx) error {
	rec, err := h.maintenance.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, rec)
}

func (h *Handler) ListMaintenanceRecords(c *fiber.Ctx) error {
	filter := mreq.ListMaintenanceFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("asset_id"); v != "" {
		filter.AssetID = &v
	}
	if v := c.Query("maintenance_status"); v != "" {
		filter.MaintenanceStatus = &v
	}
	records, err := h.maintenance.List(c.Context(), tenantID(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, records, fiber.Map{
		"page": filter.Page, "per_page": filter.PerPage,
	})
}

func (h *Handler) UpdateMaintenanceStatus(c *fiber.Ctx) error {
	var req mreq.UpdateMaintenanceStatus
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	rec, err := h.maintenance.UpdateStatus(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, rec)
}

func (h *Handler) DeleteMaintenanceRecord(c *fiber.Ctx) error {
	if err := h.maintenance.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Inspections ----

func (h *Handler) CreateInspection(c *fiber.Ctx) error {
	var req mreq.CreateInspection
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	i, err := h.inspections.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, i)
}

func (h *Handler) GetInspection(c *fiber.Ctx) error {
	i, err := h.inspections.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, i)
}

func (h *Handler) ListInspections(c *fiber.Ctx) error {
	filter := mreq.ListInspectionsFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("asset_id"); v != "" {
		filter.AssetID = &v
	}
	if v := c.Query("inspection_type"); v != "" {
		filter.InspectionType = &v
	}
	if v := c.Query("result"); v != "" {
		filter.Result = &v
	}
	items, err := h.inspections.List(c.Context(), tenantID(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, items, fiber.Map{
		"page": filter.Page, "per_page": filter.PerPage,
	})
}

// ---- Damage Reports ----

func (h *Handler) CreateDamageReport(c *fiber.Ctx) error {
	var req mreq.CreateDamageReport
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	rep, err := h.damageReports.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, rep)
}

func (h *Handler) GetDamageReport(c *fiber.Ctx) error {
	rep, err := h.damageReports.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, rep)
}

func (h *Handler) ListDamageReports(c *fiber.Ctx) error {
	filter := mreq.ListDamageReportsFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("asset_id"); v != "" {
		filter.AssetID = &v
	}
	if v := c.Query("report_status"); v != "" {
		filter.ReportStatus = &v
	}
	if v := c.Query("severity"); v != "" {
		filter.Severity = &v
	}
	reports, err := h.damageReports.List(c.Context(), tenantID(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, reports, fiber.Map{
		"page": filter.Page, "per_page": filter.PerPage,
	})
}

func (h *Handler) UpdateDamageReportStatus(c *fiber.Ctx) error {
	var req mreq.UpdateDamageReportStatus
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	rep, err := h.damageReports.UpdateStatus(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, rep)
}
