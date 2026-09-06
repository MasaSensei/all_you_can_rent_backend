package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/auth/dto/request"
	"rentos-backend/internal/modules/auth/service"
	"rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

type Handler struct {
	svc      service.AuthService
	validate *validator.Validate
}

func New(svc service.AuthService, v *validator.Validate) *Handler {
	return &Handler{svc: svc, validate: v}
}

func tenantID(c *fiber.Ctx) string {
	if id, ok := c.Locals("tenant_id").(string); ok && id != "" {
		return id
	}
	return c.Get("X-Tenant-ID")
}

func userID(c *fiber.Ctx) string {
	id, _ := c.Locals("user_id").(string)
	return id
}

// POST /auth/login
func (h *Handler) Login(c *fiber.Ctx) error {
	var req request.Login
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, "invalid body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, fmt.Sprint(errs)))
	}
	tokens, err := h.svc.Login(
		c.Context(), req,
		c.Get("User-Agent"), c.IP(),
	)
	if err != nil {
		return response.FromError(c, err)
	}
	return response.Success(c, tokens)
}

// POST /auth/refresh
func (h *Handler) Refresh(c *fiber.Ctx) error {
	var req request.Refresh
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, "invalid body"))
	}
	tokens, err := h.svc.Refresh(c.Context(), req)
	if err != nil {
		return response.FromError(c, err)
	}
	return response.Success(c, tokens)
}

// POST /auth/logout
func (h *Handler) Logout(c *fiber.Ctx) error {
	var req request.Logout
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, "invalid body"))
	}
	if err := h.svc.Logout(c.Context(), req); err != nil {
		return response.FromError(c, err)
	}
	return response.NoContent(c)
}

// GET /auth/me
func (h *Handler) Me(c *fiber.Ctx) error {
	info, err := h.svc.Me(c.Context(), userID(c), tenantID(c))
	if err != nil {
		return response.FromError(c, err)
	}
	return response.Success(c, info)
}

// POST /auth/forgot-password
func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var req request.ForgotPassword
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, "invalid body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, fmt.Sprint(errs)))
	}
	// Always return success to avoid email enumeration
	_ = h.svc.ForgotPassword(c.Context(), tenantID(c), req)
	return response.Success(c, fiber.Map{
		"message": "Jika email terdaftar, instruksi reset password telah dikirim",
	})
}

// POST /auth/reset-password
func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	var req request.ResetPassword
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, "invalid body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, fmt.Sprint(errs)))
	}
	if err := h.svc.ResetPassword(c.Context(), req); err != nil {
		return response.FromError(c, err)
	}
	return response.Success(c, fiber.Map{"message": "Password berhasil direset"})
}

// POST /auth/change-password  (requires auth)
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	var req request.ChangePassword
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, "invalid body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return response.Error(c, response.NewAppError(response.CodeValidation, fmt.Sprint(errs)))
	}
	if err := h.svc.ChangePassword(c.Context(), userID(c), tenantID(c), req); err != nil {
		return response.FromError(c, err)
	}
	return response.NoContent(c)
}
