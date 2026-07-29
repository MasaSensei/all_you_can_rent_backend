package handler

import (
	"github.com/gofiber/fiber/v2"

	ireq "rentos-backend/internal/modules/integration/dto/request"
	"rentos-backend/internal/modules/integration/service"
	apiresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

const (
	ctxKeyTenantID = "tenant_id"
	ctxKeyUserID   = "user_id"
)

// Handler groups all integration HTTP handlers.
type Handler struct {
	apiKeys  service.APIKeyService
	webhooks service.WebhookService
	validate *validator.Validate
}

func New(apiKeys service.APIKeyService, webhooks service.WebhookService, v *validator.Validate) *Handler {
	return &Handler{apiKeys: apiKeys, webhooks: webhooks, validate: v}
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

// ---- API Keys ----

func (h *Handler) CreateAPIKey(c *fiber.Ctx) error {
	var req ireq.CreateAPIKey
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	created, err := h.apiKeys.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	// 201 — raw key visible only in this response, never again.
	return apiresponse.Created(c, created)
}

func (h *Handler) GetAPIKey(c *fiber.Ctx) error {
	k, err := h.apiKeys.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, k)
}

func (h *Handler) ListAPIKeys(c *fiber.Ctx) error {
	keys, err := h.apiKeys.List(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, keys)
}

func (h *Handler) RevokeAPIKey(c *fiber.Ctx) error {
	if err := h.apiKeys.Revoke(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Webhooks ----

func (h *Handler) CreateWebhook(c *fiber.Ctx) error {
	var req ireq.CreateWebhook
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	w, err := h.webhooks.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, w)
}

func (h *Handler) GetWebhook(c *fiber.Ctx) error {
	w, err := h.webhooks.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, w)
}

func (h *Handler) ListWebhooks(c *fiber.Ctx) error {
	webhooks, err := h.webhooks.List(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, webhooks)
}

func (h *Handler) UpdateWebhook(c *fiber.Ctx) error {
	var req ireq.UpdateWebhook
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	w, err := h.webhooks.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, w)
}

func (h *Handler) DeleteWebhook(c *fiber.Ctx) error {
	if err := h.webhooks.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

func (h *Handler) ListWebhookLogs(c *fiber.Ctx) error {
	logs, err := h.webhooks.ListLogs(c.Context(), c.Params("id"), tenantID(c),
		c.QueryInt("page", 1), c.QueryInt("per_page", 20),
	)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, logs)
}
