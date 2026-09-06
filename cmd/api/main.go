package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/config"
	appmw "rentos-backend/internal/middleware"
	authmodule "rentos-backend/internal/modules/auth"
	bookingmodule "rentos-backend/internal/modules/booking"
	cmsmodule "rentos-backend/internal/modules/cms"
	coremodule "rentos-backend/internal/modules/core"
	crmmodule "rentos-backend/internal/modules/crm"
	financemodule "rentos-backend/internal/modules/finance"
	integrationmodule "rentos-backend/internal/modules/integration"
	inventorymodule "rentos-backend/internal/modules/inventory"
	maintenancemodule "rentos-backend/internal/modules/maintenance"
	notificationmodule "rentos-backend/internal/modules/notification"
	pricingmodule "rentos-backend/internal/modules/pricing"
	rbacmodule "rentos-backend/internal/modules/rbac"
	reportsmodule "rentos-backend/internal/modules/reports"
	samodule "rentos-backend/internal/modules/superadmin"
	"rentos-backend/pkg/response"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	container, err := bootstrap.New(cfg)
	if err != nil {
		panic(err)
	}
	defer container.Close()

	app := fiber.New(fiber.Config{
		AppName:      "RentOS API",
		BodyLimit:    10 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error { return response.FromError(c, err) },
	})

	app.Use(appmw.RequestID())
	app.Use(appmw.SecurityHeaders())
	app.Use(appmw.NoCache())
	app.Use(appmw.Recover(container.Logger))
	app.Use(appmw.RequestLogger(container.Logger))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Tenant-ID,X-API-Key,X-Request-ID",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))
	app.Use(compress.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.Map{"status": "ok", "checks": container.HealthCheck(c.Context())})
	})

	// ---- Modules ----
	core := coremodule.New(container)
	auth := authmodule.New(container)
	sa := samodule.New(container)
	rbac := rbacmodule.New(container)
	inventory := inventorymodule.New(container)
	pricing := pricingmodule.New(container)
	finance := financemodule.New(container)
	crm := crmmodule.New(container)
	maintenance := maintenancemodule.New(container)
	notif := notificationmodule.New(container)
	cms := cmsmodule.New(container)
	reports := reportsmodule.New(container)
	integration := integrationmodule.New(container)
	booking := bookingmodule.New(container, inventory.AssetService(), pricing.PricingQuoter())
	v1 := app.Group("/api/v1")

	// ---- Public (no JWT, with rate limit) ----
	public := v1.Group("", appmw.RateLimitDefault())
	core.RegisterPublic(public)
	auth.RegisterPublic(public)

	// Tenant registration (stricter rate limit)
	v1.Post("/auth/register", appmw.RateLimitStrict(), sa.RegisterPublic)

	// Public plan listing
	sa.RegisterPublicPlans(v1)

	// ---- Protected: JWT required ----
	protected := v1.Group("",
		appmw.TenantResolver(),
		appmw.Auth(container.JWT),
		appmw.RateLimitDefault(),
	)
	auth.RegisterProtected(protected)
	rbac.RegisterRoutes(protected)
	inventory.RegisterRoutes(protected)
	pricing.RegisterRoutes(protected)
	finance.RegisterRoutes(protected)
	crm.RegisterRoutes(protected)
	maintenance.RegisterRoutes(protected)
	notif.RegisterRoutes(protected)
	booking.RegisterRoutes(protected)
	cms.RegisterRoutes(protected)
	reports.RegisterRoutes(protected)
	integration.RegisterRoutes(protected)

	// ---- Super admin ----
	sa.RegisterSuperAdmin(v1.Group("/super-admin", appmw.RateLimitDefault()))

	// ---- External: API key ----
	reports.RegisterRoutes(v1.Group("/external",
		appmw.APIKeyAuth(integration.APIKeyService()),
		appmw.RateLimitDefault(),
	))

	_ = notif.NotificationService()
	_ = integration.WebhookService()
	_ = maintenance.MaintenanceService()
	_ = reports.ReportService()

	go func() {
		if err := app.Listen(":" + cfg.App.Port); err != nil {
			container.Logger.Fatal().Err(err).Msg("server failed")
		}
	}()
	container.Logger.Info().Str("port", cfg.App.Port).Msg("RentOS API started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(ctx)
}
