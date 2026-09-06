package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/superadmin/dto/request"
	"rentos-backend/internal/modules/superadmin/service"
	"rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

type Handler struct {
	svc      service.SuperAdminService
	validate *validator.Validate
}

func New(svc service.SuperAdminService, v *validator.Validate) *Handler {
	return &Handler{svc: svc, validate: v}
}

// POST /super-admin/login
func (h *Handler) Login(c *fiber.Ctx) error {
	var req request.SuperAdminLogin
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, "invalid body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, fmt.Sprintf("%v", errs)))
	}
	out, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return response.FromError(c, err)
	}
	return response.Success(c, out)
}

// POST /auth/register  (public)
func (h *Handler) RegisterTenant(c *fiber.Ctx) error {
	var req request.RegisterTenant
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, "invalid body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, fmt.Sprintf("%v", errs)))
	}
	out, err := h.svc.RegisterTenant(c.Context(), req)
	if err != nil {
		return response.FromError(c, err)
	}
	return response.Created(c, out)
}

// GET /super-admin/tenants
func (h *Handler) ListTenants(c *fiber.Ctx) error {
	filter := request.ListTenantsFilter{
		Search:             c.Query("search"),
		SubscriptionStatus: c.Query("subscription_status"),
		Status:             c.Query("status"),
		Page:               c.QueryInt("page", 1),
		PerPage:            c.QueryInt("per_page", 20),
	}
	out, err := h.svc.ListTenants(c.Context(), filter)
	if err != nil {
		return response.FromError(c, err)
	}
	return response.Success(c, out, fiber.Map{"page": filter.Page, "per_page": filter.PerPage})
}

// GET /super-admin/tenants/:id
func (h *Handler) GetTenant(c *fiber.Ctx) error {
	out, err := h.svc.GetTenant(c.Context(), c.Params("id"))
	if err != nil {
		return response.FromError(c, err)
	}
	return response.Success(c, out)
}

// PATCH /super-admin/tenants/:id/status
func (h *Handler) UpdateTenantStatus(c *fiber.Ctx) error {
	var req request.UpdateTenantStatus
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, "invalid body"))
	}
	if err := h.svc.UpdateTenantStatus(c.Context(), c.Params("id"), req); err != nil {
		return response.FromError(c, err)
	}
	return response.NoContent(c)
}

// POST /super-admin/tenants/:id/subscriptions
func (h *Handler) AssignSubscription(c *fiber.Ctx) error {
	var req request.AssignSubscription
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, "invalid body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, fmt.Sprintf("%v", errs)))
	}
	if err := h.svc.AssignSubscription(c.Context(), c.Params("id"), req); err != nil {
		return response.FromError(c, err)
	}
	return response.NoContent(c)
}

// GET /super-admin/plans
func (h *Handler) ListPlans(c *fiber.Ctx) error {
	out, err := h.svc.ListPlans(c.Context())
	if err != nil {
		return response.FromError(c, err)
	}
	return response.Success(c, out)
}

// GET /super-admin/stats
func (h *Handler) Stats(c *fiber.Ctx) error {
	out, err := h.svc.GetStats(c.Context())
	if err != nil {
		return response.FromError(c, err)
	}
	return response.Success(c, out)
}
