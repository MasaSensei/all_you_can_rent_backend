package handler

import (
	"github.com/gofiber/fiber/v2"

	rbacdto "rentos/internal/modules/rbac/dto/request"
	rbacresponse "rentos/internal/modules/rbac/dto/response"
	"rentos/internal/modules/rbac/entity"
	"rentos/internal/modules/rbac/service"
	apiresponse "rentos/pkg/response"
	"rentos/pkg/validator"
)

const ctxKeyTenantID = "tenant_id"
const ctxKeyUserID = "user_id"

// Handler groups the RBAC module's HTTP handlers.
type Handler struct {
	roles       service.RoleService
	permissions service.PermissionService
	userRoles   service.UserRoleService
	validate    *validator.Validate
}

func New(
	roles service.RoleService,
	permissions service.PermissionService,
	userRoles service.UserRoleService,
	v *validator.Validate,
) *Handler {
	return &Handler{roles: roles, permissions: permissions, userRoles: userRoles, validate: v}
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

// ---- Roles ----

func (h *Handler) CreateRole(c *fiber.Ctx) error {
	var req rbacdto.CreateRole
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	role, err := h.roles.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, toRoleResponse(role, nil))
}

func (h *Handler) GetRole(c *fiber.Ctx) error {
	role, err := h.roles.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	perms, err := h.roles.GetRolePermissions(c.Context(), role.ID, tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toRoleResponse(role, perms))
}

func (h *Handler) ListRoles(c *fiber.Ctx) error {
	roles, err := h.roles.List(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	out := make([]rbacresponse.Role, 0, len(roles))
	for _, r := range roles {
		out = append(out, toRoleResponse(&r, nil))
	}
	return apiresponse.Success(c, out)
}

func (h *Handler) UpdateRole(c *fiber.Ctx) error {
	var req rbacdto.UpdateRole
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	role, err := h.roles.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toRoleResponse(role, nil))
}

func (h *Handler) DeleteRole(c *fiber.Ctx) error {
	if err := h.roles.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

func (h *Handler) SyncPermissions(c *fiber.Ctx) error {
	var req rbacdto.SyncPermissions
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	perms, err := h.roles.SyncPermissions(c.Context(), c.Params("id"), tenantID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	out := make([]rbacresponse.Permission, 0, len(perms))
	for _, p := range perms {
		out = append(out, toPermissionResponse(&p))
	}
	return apiresponse.Success(c, out)
}

// ---- Permissions ----

func (h *Handler) ListPermissions(c *fiber.Ctx) error {
	perms, err := h.permissions.List(c.Context())
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	out := make([]rbacresponse.Permission, 0, len(perms))
	for _, p := range perms {
		out = append(out, toPermissionResponse(&p))
	}
	return apiresponse.Success(c, out)
}

// ---- User Roles ----

func (h *Handler) AssignRole(c *fiber.Ctx) error {
	var req rbacdto.AssignRole
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	if err := h.userRoles.Assign(c.Context(), c.Params("user_id"), tenantID(c), req); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, fiber.Map{"message": "role assigned"})
}

func (h *Handler) RevokeRole(c *fiber.Ctx) error {
	if err := h.userRoles.Revoke(c.Context(), c.Params("user_id"), c.Params("role_id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

func (h *Handler) ListUserRoles(c *fiber.Ctx) error {
	roles, err := h.userRoles.ListForUser(c.Context(), c.Params("user_id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	out := make([]rbacresponse.Role, 0, len(roles))
	for _, r := range roles {
		out = append(out, toRoleResponse(&r, nil))
	}
	return apiresponse.Success(c, out)
}

// ---- mapping ----

func toRoleResponse(r *entity.Role, perms []entity.Permission) rbacresponse.Role {
	out := rbacresponse.Role{
		ID:          r.ID,
		TenantID:    r.TenantID,
		Name:        r.Name,
		Description: r.Description,
		IsSystem:    r.IsSystem,
		Status:      r.Status,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
	if perms != nil {
		out.Permissions = make([]rbacresponse.Permission, 0, len(perms))
		for _, p := range perms {
			out.Permissions = append(out.Permissions, toPermissionResponse(&p))
		}
	}
	return out
}

func toPermissionResponse(p *entity.Permission) rbacresponse.Permission {
	return rbacresponse.Permission{
		ID:     p.ID,
		Name:   p.Name,
		Module: p.Module,
		Action: p.Action,
	}
}
