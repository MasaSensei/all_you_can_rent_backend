package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"rentos/internal/bootstrap"
	"rentos/internal/config"
	appmiddleware "rentos/internal/middleware"
	authmodule "rentos/internal/modules/auth"
	bookingmodule "rentos/internal/modules/booking"
	coremodule "rentos/internal/modules/core"
	inventorymodule "rentos/internal/modules/inventory"
	rbacmodule "rentos/internal/modules/rbac"
	"rentos/pkg/response"
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
		ErrorHandler: func(c *fiber.Ctx, err error) error { return response.FromError(c, err) },
	})

	// ---- Global middleware ----
	app.Use(appmiddleware.RequestID())
	app.Use(appmiddleware.Recover(container.Logger))
	app.Use(appmiddleware.RequestLogger(container.Logger))
	app.Use(cors.New())

	// ---- Health ----
	app.Get("/health", func(c *fiber.Ctx) error {
		return response.Success(c, fiber.Map{
			"status": "ok", "checks": container.HealthCheck(c.Context()),
		})
	})

	// ---- Build modules (dependency order matters) ----
	core := coremodule.New(container)
	auth := authmodule.New(container)
	rbac := rbacmodule.New(container)
	inventory := inventorymodule.New(container)

	// Booking receives inventory.AssetService via the InventoryChecker interface —
	// no import of the inventory postgres package from booking.
	booking := bookingmodule.New(container, inventory.AssetService())

	// ---- API v1 ----
	v1 := app.Group("/api/v1")
	v1.Use(appmiddleware.TenantResolver())

	// Public routes (no JWT required)
	core.RegisterRoutes(v1)
	auth.RegisterRoutes(v1)

	// Protected routes (JWT required)
	protected := v1.Group("", appmiddleware.Auth(container.JWT))
	rbac.RegisterRoutes(protected)
	inventory.RegisterRoutes(protected)
	booking.RegisterRoutes(protected)

	// Future modules wired here:
	//   pricing := pricingmodule.New(container)
	//   pricing.RegisterRoutes(protected)
	//   finance := financemodule.New(container)
	//   finance.RegisterRoutes(protected)

	// ---- Start ----
	go func() {
		if err := app.Listen(":" + cfg.App.Port); err != nil {
			container.Logger.Fatal().Err(err).Msg("server failed to start")
		}
	}()
	container.Logger.Info().
		Str("port", cfg.App.Port).
		Str("env", cfg.App.Env).
		Msg("RentOS API started")

	waitForShutdown(app, container)
}

func waitForShutdown(app *fiber.App, c *bootstrap.Container) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c.Logger.Info().Msg("shutting down gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		c.Logger.Error().Err(err).Msg("forced shutdown")
	}
}
