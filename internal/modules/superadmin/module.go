package superadmin

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	authpostgres "rentos-backend/internal/modules/auth/repository/postgres"
	"rentos-backend/internal/modules/superadmin/handler"
	sapostgres "rentos-backend/internal/modules/superadmin/repository/postgres"
	"rentos-backend/internal/modules/superadmin/routes"
	"rentos-backend/internal/modules/superadmin/service"
	jwtservice "rentos-backend/pkg/jwt"
)

type accessTokenIssuerAdapter struct {
	service *jwtservice.Service
}

func (a accessTokenIssuerAdapter) IssueAccess(subject, email, username string) (string, error) {
	return a.service.IssueAccess(subject, email, username, nil)
}

type Module struct {
	handler *handler.Handler
	svc     service.SuperAdminService
}

func New(c *bootstrap.Container) *Module {
	tenantRepo := sapostgres.NewTenantRepository(
		query("create_tenant.sql"), query("find_tenant_by_id.sql"),
		query("find_tenant_by_slug.sql"), query("slug_exists.sql"),
		query("list_tenants.sql"), query("update_tenant_status.sql"),
		query("update_tenant_subscription_status.sql"), query("platform_stats.sql"),
	)
	planRepo := sapostgres.NewSubscriptionPlanRepository(
		query("find_plan_by_slug.sql"), query("list_active_plans.sql"),
	)
	subsRepo := sapostgres.NewTenantSubscriptionRepository(
		query("create_subscription.sql"), query("find_active_subscription.sql"),
	)
	saRepo := sapostgres.NewSuperAdminRepository(
		query("find_super_admin_by_email.sql"), query("update_super_admin_last_login.sql"),
	)
	userRepo := authpostgres.NewUserRepository(
		query("create_user.sql"),
		query("find_user_by_id.sql"),
		query("find_user_by_email.sql"),
		query("find_user_by_username.sql"),
		query("update_password.sql"),
	)

	svc := service.NewSuperAdminService(c.DB, tenantRepo, planRepo, subsRepo, saRepo, userRepo, accessTokenIssuerAdapter{service: c.JWT})
	h := handler.New(svc, c.Validator)
	return &Module{handler: h, svc: svc}
}

// RegisterPublic registers the /auth/register route on the provided router.
func (m *Module) RegisterPublic(c *fiber.Ctx) error {
	return m.handler.RegisterTenant(c)
}

// RegisterPublicPlans registers the public /plans listing.
func (m *Module) RegisterPublicPlans(router fiber.Router) {
	router.Get("/plans", m.handler.ListPlans)
}

// RegisterSuperAdmin mounts the /super-admin/* routes.
func (m *Module) RegisterSuperAdmin(router fiber.Router) {
	routes.RegisterSuperAdmin(router, m.handler)
}

func (m *Module) SuperAdminService() service.SuperAdminService { return m.svc }
