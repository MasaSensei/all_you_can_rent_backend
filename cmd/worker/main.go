package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/config"
	integrationmodule  "rentos-backend/internal/modules/integration"
	maintenancemodule  "rentos-backend/internal/modules/maintenance"
	notificationmodule "rentos-backend/internal/modules/notification"
	reportsmodule      "rentos-backend/internal/modules/reports"
	workerMaintenance  "rentos-backend/internal/worker/maintenance"
	workerNotification "rentos-backend/internal/worker/notification"
	workerReport       "rentos-backend/internal/worker/report"
	workerWebhook      "rentos-backend/internal/worker/webhook"
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

	log := container.Logger

	// ---- Build modules needed by workers ----
	maintenance  := maintenancemodule.New(container)
	reports      := reportsmodule.New(container)
	notification := notificationmodule.New(container)
	integration  := integrationmodule.New(container)

	// ---- Instantiate workers ----
	maintScheduler := workerMaintenance.NewScheduler(
		maintenance.MaintenanceService(),
		5*time.Minute,
		log.With().Str("worker", "maintenance_scheduler").Logger(),
	)

	reportGen := workerReport.NewGenerator(
		reports.ReportService(),
		30*time.Second,
		log.With().Str("worker", "report_generator").Logger(),
	)

	webhookDispatcher := workerWebhook.NewDispatcher(
		// Webhook + log repos are accessed through the integration module's
		// internal wiring; for the dispatcher we pass the DB directly.
		// In a more complete setup, expose repos from integration.Module.
		nil, nil, container.DB,
		10*time.Second,
		log.With().Str("worker", "webhook_dispatcher").Logger(),
	)
	_ = integration.WebhookService() // referenced so import is used

	notifWorker := workerNotification.NewDeliveryWorker(
		notification.LogService(),
		15*time.Second,
		log.With().Str("worker", "notification_delivery").Logger(),
	)

	// ---- Run all workers concurrently ----
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	run := func(name string, fn func(context.Context)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info().Str("worker", name).Msg("starting")
			fn(ctx)
			log.Info().Str("worker", name).Msg("stopped")
		}()
	}

	run("maintenance_scheduler",  maintScheduler.Run)
	run("report_generator",       reportGen.Run)
	run("webhook_dispatcher",     webhookDispatcher.Run)
	run("notification_delivery",  notifWorker.Run)

	log.Info().Msg("all workers running")

	// ---- Graceful shutdown ----
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutdown signal received — stopping workers")
	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("all workers stopped cleanly")
	case <-time.After(30 * time.Second):
		log.Warn().Msg("worker shutdown timed out after 30s")
	}
}
