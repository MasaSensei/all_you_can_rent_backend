// Package logger provides the application's structured logger. Every other
// package depends on this instead of importing zerolog directly, so the
// underlying logging library can be swapped later without touching call
// sites.
package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type ctxKey struct{}

// Config controls log level and output formatting.
type Config struct {
	Level      string // debug, info, warn, error
	PrettyText bool   // human-readable output for local dev; JSON otherwise
}

// New builds the base application logger.
func New(cfg Config) zerolog.Logger {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.TimeFieldFormat = time.RFC3339

	if cfg.PrettyText {
		return zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
			Level(level).With().Timestamp().Logger()
	}
	return zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()
}

// WithContext stores a request-scoped logger (already enriched with
// fields such as request_id/tenant_id/user_id) on the context.
func WithContext(ctx context.Context, l zerolog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext returns the request-scoped logger, falling back to a base
// logger if none was attached to the context.
func FromContext(ctx context.Context) zerolog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(zerolog.Logger); ok {
		return l
	}
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
