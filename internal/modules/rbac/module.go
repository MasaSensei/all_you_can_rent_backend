package rbac

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/rbac/handler"
	"rentos-backend/internal/modules/rbac/repository/postgres"
	"rentos-backend/internal/modules/rbac/routes"
	"rentos-backend/internal/modules/rbac/service"
)

// Module holds the RBAC module's wired handler.
type Module struct {
	handler      *handler.Handler
	userRoleSvc  service.UserRoleService
}

// New builds the RBAC module: repositories → services → handler.
func New(c *bootstrap.Container) *Module {
	roleRepo := postgres.NewRoleRepository(
		query("create_role.sql"),
		query("find_role_by_id.sql"),
		query("list_roles.sql"),
		query("update_role.sql"),
		query("delete_role.sql"),
	)
	permRepo := postgres.NewPermissionRepository(
		query("list_permissions.sql"),
		query("find_permissions_by_ids.sql"),
		query("get_role_permissions.sql"),
	)
	rolePermRepo := postgres.NewRolePermissionRepository(
		query("delete_role_permissions.sql"),
		query("create_role_permission.sql"),
	)
	userRoleRepo := postgres.NewUserRoleRepository(
		query("assign_user_role.sql"),
		query("revoke_user_role.sql"),
		query("list_user_roles.sql"),
		query("get_user_permissions.sql"),
	)

	roleSvc := service.NewRoleService(c.DB, roleRepo, permRepo, rolePermRepo)
	permSvc := service.NewPermissionService(c.DB, permRepo)
	userRoleSvc := service.NewUserRoleService(c.DB, userRoleRepo, roleRepo)

	h := handler.New(roleSvc, permSvc, userRoleSvc, c.Validator)
	return &Module{handler: h, userRoleSvc: userRoleSvc}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}

// UserRoleService exposes the UserRoleService so RBACMiddleware can
// resolve permissions without importing the rbac package itself —
// middleware depends on the interface, not this concrete module.
func (m *Module) UserRoleService() service.UserRoleService {
	return m.userRoleSvc
}
