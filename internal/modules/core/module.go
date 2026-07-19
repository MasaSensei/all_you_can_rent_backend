// Package core wires the core module (tenants, settings, system_settings,
// audit_logs) together from the shared Container and exposes a single
// RegisterRoutes entrypoint. This is the only file other packages
// (internal/bootstrap) need to import to mount the module.
package core

import (
	"github.com/gofiber/fiber/v2"

	"rentos/internal/bootstrap"
	"rentos/internal/modules/core/handler"
	"rentos/internal/modules/core/repository/postgres"
	"rentos/internal/modules/core/routes"
	"rentos/internal/modules/core/service"
)

// Module holds the core module's wired handler, ready to register routes.
type Module struct {
	handler *handler.Handler
}

// New builds the core module: repositories -> services -> handler.
// Query strings are read once here from the embedded sql/ folder
// (see queries.go) and injected into each repository constructor.
// The repository/postgres package is a pure SQL-execution layer with no
// embedding concerns of its own.
func New(c *bootstrap.Container) *Module {
	tenantRepo := postgres.NewTenantRepository(
		query("create_tenant.sql"),
		query("find_tenant_by_id.sql"),
		query("find_tenant_by_slug.sql"),
		query("update_tenant_status.sql"),
	)
	settingRepo := postgres.NewSettingRepository(
		query("upsert_setting.sql"),
		query("find_setting.sql"),
		query("list_settings.sql"),
	)
	systemSettingRepo := postgres.NewSystemSettingRepository(
		query("upsert_system_setting.sql"),
		query("find_system_setting.sql"),
	)
	auditRepo := postgres.NewAuditLogRepository(
		query("create_audit_log.sql"),
		query("list_audit_logs.sql"),
	)

	tenantSvc := service.NewTenantService(c.DB, tenantRepo)
	settingSvc := service.NewSettingService(c.DB, settingRepo, systemSettingRepo)
	auditSvc := service.NewAuditService(c.DB, auditRepo)

	h := handler.New(tenantSvc, settingSvc, auditSvc, c.Validator)
	return &Module{handler: h}
}

// RegisterRoutes mounts the module's routes onto router (expected to
// already be scoped to /api/v1).
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}
