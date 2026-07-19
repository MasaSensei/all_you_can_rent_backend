// Package bootstrap wires shared infrastructure together into a single
// Container. Feature modules receive this Container in their constructor
// and pull out only what they need.
package bootstrap

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"rentos/internal/config"
	"rentos/pkg/cache"
	"rentos/pkg/database"
	"rentos/pkg/jwt"
	"rentos/pkg/logger"
	"rentos/pkg/validator"
)

// Container holds every dependency a feature module needs.
type Container struct {
	Config    *config.Config
	Logger    zerolog.Logger
	DB        *sqlx.DB
	Redis     *redis.Client
	JWT       *jwt.Service
	Validator *validator.Validate
}

// New builds and connects every infrastructure dependency.
func New(cfg *config.Config) (*Container, error) {
	log := logger.New(logger.Config{
		Level:      cfg.App.LogLevel,
		PrettyText: cfg.App.Env == "local",
	})

	db, err := database.New(database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		Name:            cfg.Database.Name,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		return nil, fmt.Errorf("bootstrap: database: %w", err)
	}

	redisClient := cache.New(cache.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	jwtSvc := jwt.New(jwt.Config{
		AccessSecret:  cfg.JWT.AccessSecret,
		RefreshSecret: cfg.JWT.RefreshSecret,
		AccessTTL:     cfg.JWT.AccessTTL,
		RefreshTTL:    cfg.JWT.RefreshTTL,
	})

	return &Container{
		Config:    cfg,
		Logger:    log,
		DB:        db,
		Redis:     redisClient,
		JWT:       jwtSvc,
		Validator: validator.New(),
	}, nil
}

// HealthCheck verifies every infrastructure dependency is reachable.
func (c *Container) HealthCheck(ctx context.Context) map[string]string {
	status := map[string]string{"database": "ok", "redis": "ok"}
	if err := database.Health(ctx, c.DB); err != nil {
		status["database"] = err.Error()
	}
	if err := cache.Health(ctx, c.Redis); err != nil {
		status["redis"] = err.Error()
	}
	return status
}

// Close releases all infrastructure connections on graceful shutdown.
func (c *Container) Close() error {
	if err := c.DB.Close(); err != nil {
		return err
	}
	return c.Redis.Close()
}
