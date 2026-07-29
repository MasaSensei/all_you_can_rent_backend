package handler

import (
	"github.com/gofiber/fiber/v2"

	nreq "rentos-backend/internal/modules/notification/dto/request"
	"rentos-backend/internal/modules/notification/service"
	apiresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

const (
	ctxKeyTenantID = "tenant_id"
	ctxKeyUserID   = "user_id"
)

// Handler groups all notification HTTP handlers.
type Handler struct {
	templates service.TemplateService
	notifs    service.NotificationService
	logs      service.LogService
	validate  *validator.Validate
}

func New(
	templates service.TemplateService,
	notifs service.NotificationService,
	logs service.LogService,
	v *validator.Validate,
) *Handler {
	return &Handler{templates: templates, notifs: notifs, logs: logs, validate: v}
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

// ---- Notification Templates ----

func (h *Handler) CreateTemplate(c *fiber.Ctx) error {
	var req nreq.CreateTemplate
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	t, err := h.templates.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, t)
}

func (h *Handler) GetTemplate(c *fiber.Ctx) error {
	t, err := h.templates.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, t)
}

func (h *Handler) ListTemplates(c *fiber.Ctx) error {
	templates, err := h.templates.List(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, templates)
}

func (h *Handler) UpdateTemplate(c *fiber.Ctx) error {
	var req nreq.UpdateTemplate
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	t, err := h.templates.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, t)
}

func (h *Handler) DeleteTemplate(c *fiber.Ctx) error {
	if err := h.templates.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- In-App Notifications ----

func (h *Handler) CreateInApp(c *fiber.Ctx) error {
	var req nreq.CreateInAppNotification
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	n, err := h.notifs.CreateInApp(c.Context(), tenantID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, n)
}

func (h *Handler) ListNotifications(c *fiber.Ctx) error {
	filter := nreq.ListNotificationsFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("is_read"); v != "" {
		b := v == "true"
		filter.IsRead = &b
	}

	// Recipient resolved from auth context.
	uid := userID(c)
	var uIDPtr *string
	if uid != "" {
		uIDPtr = &uid
	}

	notifs, err := h.notifs.List(c.Context(), tenantID(c), filter, uIDPtr, nil)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, notifs, fiber.Map{
		"page": filter.Page, "per_page": filter.PerPage,
	})
}

func (h *Handler) MarkRead(c *fiber.Ctx) error {
	if err := h.notifs.MarkRead(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

func (h *Handler) MarkAllRead(c *fiber.Ctx) error {
	uid := userID(c)
	var uIDPtr *string
	if uid != "" {
		uIDPtr = &uid
	}
	if err := h.notifs.MarkAllRead(c.Context(), tenantID(c), uIDPtr, nil); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Notification Logs ----

func (h *Handler) ListLogs(c *fiber.Ctx) error {
	logs, err := h.logs.List(c.Context(), tenantID(c),
		c.QueryInt("page", 1), c.QueryInt("per_page", 20),
	)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, logs)
}
