// Package handler adapts HTTP requests to the core module's services.
// Handlers only parse/validate input and shape output — all business
// logic lives in service.
package handler

import (
	"github.com/gofiber/fiber/v2"

	"rentos/internal/modules/core/dto/request"
	"rentos/internal/modules/core/dto/response"
	"rentos/internal/modules/core/entity"
	"rentos/internal/modules/core/service"
	apiresponse "rentos/pkg/response"
	"rentos/pkg/validator"
)

// HeaderTenantID is a temporary stand-in for tenant resolution until the
// TenantResolver middleware (Phase 2/3, tenant module) is in place. Every
// handler below reads the tenant ID exclusively through this helper, so
// swapping the resolution source later means changing one function, not
// every handler.
const HeaderTenantID = "X-Tenant-ID"

// Handler groups the core module's HTTP handlers.
type Handler struct {
	tenants  service.TenantService
	settings service.SettingService
	audit    service.AuditService
	validate *validator.Validate
}

// New builds the core module's Handler.
func New(tenants service.TenantService, settings service.SettingService, audit service.AuditService, v *validator.Validate) *Handler {
	return &Handler{tenants: tenants, settings: settings, audit: audit, validate: v}
}

func tenantIDFromRequest(c *fiber.Ctx) string {
	return c.Get(HeaderTenantID)
}

// ---- Tenants ----

func (h *Handler) CreateTenant(c *fiber.Ctx) error {
	var req request.CreateTenant
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	t, err := h.tenants.CreateTenant(c.Context(), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, toTenantResponse(t))
}

func (h *Handler) GetTenant(c *fiber.Ctx) error {
	t, err := h.tenants.GetTenantByID(c.Context(), c.Params("id"))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toTenantResponse(t))
}

// ---- Settings ----

func (h *Handler) ListSettings(c *fiber.Ctx) error {
	settings, err := h.settings.ListSettings(c.Context(), tenantIDFromRequest(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}

	out := make([]response.Setting, 0, len(settings))
	for _, s := range settings {
		out = append(out, toSettingResponse(&s))
	}
	return apiresponse.Success(c, out)
}

func (h *Handler) GetSetting(c *fiber.Ctx) error {
	s, err := h.settings.GetSetting(c.Context(), tenantIDFromRequest(c), c.Params("key"))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toSettingResponse(s))
}

func (h *Handler) UpsertSetting(c *fiber.Ctx) error {
	var req request.UpsertSetting
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	s, err := h.settings.UpsertSetting(c.Context(), tenantIDFromRequest(c), req, nil)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toSettingResponse(s))
}

// ---- System Settings (global, intended for super-admin use once RBAC exists) ----

func (h *Handler) GetSystemSetting(c *fiber.Ctx) error {
	s, err := h.settings.GetSystemSetting(c.Context(), c.Params("key"))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toSystemSettingResponse(s))
}

func (h *Handler) UpsertSystemSetting(c *fiber.Ctx) error {
	var req request.UpsertSystemSetting
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	s, err := h.settings.UpsertSystemSetting(c.Context(), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toSystemSettingResponse(s))
}

// ---- Audit Logs (read-only via API; writes happen internally from other services) ----

func (h *Handler) ListAuditLogs(c *fiber.Ctx) error {
	filter := request.ListAuditLogsFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("entity_type"); v != "" {
		filter.EntityType = &v
	}
	if v := c.Query("entity_id"); v != "" {
		filter.EntityID = &v
	}
	if v := c.Query("action"); v != "" {
		filter.Action = &v
	}

	logs, err := h.audit.List(c.Context(), tenantIDFromRequest(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}

	out := make([]response.AuditLog, 0, len(logs))
	for _, l := range logs {
		out = append(out, toAuditLogResponse(&l))
	}
	return apiresponse.Success(c, out, fiber.Map{"page": filter.Page, "per_page": filter.PerPage})
}

// ---- entity -> dto/response mapping ----

func toTenantResponse(t *entity.Tenant) response.Tenant {
	return response.Tenant{
		ID: t.ID, Name: t.Name, Slug: t.Slug, Domain: t.Domain,
		Plan: t.Plan, Status: t.Status, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func toSettingResponse(s *entity.Setting) response.Setting {
	return response.Setting{ID: s.ID, Key: s.Key, Value: s.Value, Type: s.Type, UpdatedAt: s.UpdatedAt}
}

func toSystemSettingResponse(s *entity.SystemSetting) response.SystemSetting {
	return response.SystemSetting{ID: s.ID, Key: s.Key, Value: s.Value, Type: s.Type, UpdatedAt: s.UpdatedAt}
}

func toAuditLogResponse(a *entity.AuditLog) response.AuditLog {
	return response.AuditLog{
		ID: a.ID, UserID: a.UserID, EntityType: a.EntityType,
		EntityID: a.EntityID, Action: a.Action, CreatedAt: a.CreatedAt,
	}
}
