package handler

import (
	"github.com/gofiber/fiber/v2"

	authdto "rentos-backend/internal/modules/auth/dto/request"
	authresponse "rentos-backend/internal/modules/auth/dto/response"
	"rentos-backend/internal/modules/auth/entity"
	"rentos-backend/internal/modules/auth/service"
	apiresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

// ctxKeyTenantID is the fiber.Ctx.Locals key set by TenantResolver middleware.
const ctxKeyTenantID = "tenant_id"

// ctxKeySessionID is the fiber.Ctx.Locals key set by AuthMiddleware.
const ctxKeySessionID = "session_id"

// Handler groups the auth module's HTTP handlers.
type Handler struct {
	auth     service.AuthService
	users    service.UserService
	password service.PasswordService
	validate *validator.Validate
}

func New(
	auth service.AuthService,
	users service.UserService,
	pwd service.PasswordService,
	v *validator.Validate,
) *Handler {
	return &Handler{auth: auth, users: users, password: pwd, validate: v}
}

func tenantID(c *fiber.Ctx) string {
	if id, ok := c.Locals(ctxKeyTenantID).(string); ok {
		return id
	}
	// Fallback for Phase 1 compatibility until TenantResolver is wired.
	return c.Get("X-Tenant-ID")
}

// ---- Auth ----

func (h *Handler) Register(c *fiber.Ctx) error {
	var req authdto.Register
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	result, err := h.auth.Register(c.Context(), tenantID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, result)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req authdto.Login
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	result, err := h.auth.Login(c.Context(), tenantID(c), req,
		c.IP(), c.Get("User-Agent"),
	)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, result)
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	var req authdto.RefreshToken
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	result, err := h.auth.Refresh(c.Context(), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, result)
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	sessionID, _ := c.Locals(ctxKeySessionID).(string)
	if err := h.auth.Logout(c.Context(), sessionID); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var req authdto.ForgotPassword
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	if err := h.password.ForgotPassword(c.Context(), tenantID(c), req); err != nil {
		return apiresponse.FromError(c, err)
	}
	// Always return 200 to prevent email enumeration.
	return apiresponse.Success(c, fiber.Map{"message": "if the email is registered, a reset link has been sent"})
}

func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	var req authdto.ResetPassword
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	if err := h.password.ResetPassword(c.Context(), tenantID(c), req); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, fiber.Map{"message": "password reset successfully"})
}

// ---- Users ----

func (h *Handler) ListUsers(c *fiber.Ctx) error {
	users, err := h.users.List(c.Context(), tenantID(c),
		c.QueryInt("page", 1), c.QueryInt("per_page", 20),
	)
	if err != nil {
		return apiresponse.FromError(c, err)
	}

	out := make([]authresponse.User, 0, len(users))
	for _, u := range users {
		out = append(out, toUserResponse(&u))
	}
	return apiresponse.Success(c, out)
}

func (h *Handler) GetUser(c *fiber.Ctx) error {
	u, err := h.users.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toUserResponse(u))
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	var req authdto.UpdateUser
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	u, err := h.users.Update(c.Context(), c.Params("id"), tenantID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toUserResponse(u))
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	if err := h.users.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- mapping ----

func toUserResponse(u *entity.User) authresponse.User {
	return authresponse.User{
		ID:          u.ID,
		TenantID:    u.TenantID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Phone:       u.Phone,
		AvatarURL:   u.AvatarURL,
		IsActive:    u.IsActive,
		LastLoginAt: u.LastLoginAt,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
