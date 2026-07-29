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
	appmiddleware "rentos-backend/internal/middleware"
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
		AppName: "RentOS API",
		// Cap request body at 10 MB globally.
		BodyLimit:    10 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error { return response.FromError(c, err) },
	})

	// ---- Global middleware (order matters) ----
	app.Use(appmiddleware.RequestID())
	app.Use(appmiddleware.SecurityHeaders())
	app.Use(appmiddleware.NoCache())
	app.Use(appmiddleware.Recover(container.Logger))
	app.Use(appmiddleware.RequestLogger(container.Logger))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Tenant-ID, X-API-Key, X-Request-ID",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))
	app.Use(compress.New()) // gzip response compression

	// ---- Health (no auth, no rate limit) ----
	app.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status": "ok", "checks": container.HealthCheck(c.Context()),
		})
	})

	// ---- Build modules ----

	// Layer 1 — no cross-module deps
	core := coremodule.New(container)
	auth := authmodule.New(container)
	rbac := rbacmodule.New(container)
	inventory := inventorymodule.New(container)
	pricing := pricingmodule.New(container)
	finance := financemodule.New(container)
	crm := crmmodule.New(container)
	maintenance := maintenancemodule.New(container)
	notification := notificationmodule.New(container)
	cms := cmsmodule.New(container)
	reports := reportsmodule.New(container)
	integration := integrationmodule.New(container)

	// Layer 2 — receives layer 1 interfaces
	booking := bookingmodule.New(
		container,
		inventory.AssetService(),
		pricing.PricingQuoter(),
	)

	// ---- API v1 ----
	v1 := app.Group("/api/v1")
	v1.Use(appmiddleware.TenantResolver())

	// ---- Public group: no JWT, rate-limited per IP ----
	public := v1.Group("",
		appmiddleware.RateLimitDefault(),
		appmiddleware.BodySizeLimit(1*1024*1024), // 1 MB for public endpoints
	)
	core.RegisterRoutes(public)

	// Auth endpoints — stricter rate limit (brute force protection)
	authGroup := v1.Group("",
		appmiddleware.RateLimitStrict(),
		appmiddleware.BodySizeLimit(512*1024), // 512 KB
	)
	auth.RegisterRoutes(authGroup)

	// ---- Protected: JWT required ----
	protected := v1.Group("",
		appmiddleware.Auth(container.JWT),
		appmiddleware.RateLimitDefault(),
		appmiddleware.BodySizeLimit(5*1024*1024), // 5 MB
	)
	rbac.RegisterRoutes(protected)
	inventory.RegisterRoutes(protected)
	pricing.RegisterRoutes(protected)
	finance.RegisterRoutes(protected)
	crm.RegisterRoutes(protected)
	maintenance.RegisterRoutes(protected)
	notification.RegisterRoutes(protected)
	booking.RegisterRoutes(protected)
	cms.RegisterRoutes(protected)
	reports.RegisterRoutes(protected)
	integration.RegisterRoutes(protected)

	// ---- External: API key auth (for third-party integrations) ----
	external := v1.Group("/external",
		appmiddleware.APIKeyAuth(integration.APIKeyService()),
		appmiddleware.RateLimitDefault(),
	)
	// Mount a subset of routes for API key consumers (e.g. booking read, event tracking).
	reports.RegisterRoutes(external)

	// ---- Analytics: higher rate limit for event ingestion ----
	analyticsGroup := v1.Group("",
		appmiddleware.Auth(container.JWT),
		appmiddleware.RateLimitEvents(),
	)
	reports.RegisterRoutes(analyticsGroup)

	// Cross-module service references ready for injection:
	_ = notification.NotificationService() // → other modules call Send()
	_ = integration.WebhookService()       // → domain handlers call Dispatch()
	_ = maintenance.MaintenanceService()   // → worker calls ListDue()
	_ = reports.ReportService()            // → worker calls UpdateStatus()

	// ---- Start server ----
	go func() {
		if err := app.Listen(":" + cfg.App.Port); err != nil {
			container.Logger.Fatal().Err(err).Msg("server failed to start")
		}
	}()
	container.Logger.Info().
		Str("port", cfg.App.Port).
		Str("env", cfg.App.Env).
		Msg("RentOS API started — 12 modules, rate limiting active")

	waitForShutdown(app, container)
}

func waitForShutdown(app *fiber.App, c *bootstrap.Container) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c.Logger.Info().Msg("shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		c.Logger.Error().Err(err).Msg("forced shutdown")
	}
	c.Logger.Info().Msg("server stopped cleanly")
}
