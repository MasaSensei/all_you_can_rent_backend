package main

import (
	"os"
	"os/signal"
	"syscall"

	"rentos/internal/bootstrap"
	"rentos/internal/config"
)

// main starts the worker process skeleton. Job registration (email,
// whatsapp, notification, report, webhook dispatch, cron schedules) is
// added in later phases via a queue consumer wired against this same
// Container — no structural change to this bootstrap is expected then.
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

	container.Logger.Info().Str("env", cfg.App.Env).Msg("RentOS worker started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Info().Msg("worker shutting down")
}
