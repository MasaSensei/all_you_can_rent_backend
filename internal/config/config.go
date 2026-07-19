// Package config loads all application configuration from environment
// variables (optionally preloaded from a local .env file). It is the only
// package that reads os.Getenv directly — everything else receives a typed
// Config.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config is the root configuration object.
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Env      string
	Port     string
	LogLevel string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig holds signing secrets and token TTLs.
type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Env:      getEnv("APP_ENV", "local"),
			Port:     getEnv("APP_PORT", "8080"),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "rentos"),
			Password:        getEnv("DB_PASSWORD", "rentos"),
			Name:            getEnv("DB_NAME", "rentos"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MIN", 30)) * time.Minute,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", ""),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
			AccessTTL:     time.Duration(getEnvInt("JWT_ACCESS_TTL_MIN", 15)) * time.Minute,
			RefreshTTL:    time.Duration(getEnvInt("JWT_REFRESH_TTL_DAY", 30)) * 24 * time.Hour,
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.Name == "" {
		return fmt.Errorf("config: DB_NAME is required")
	}
	if c.JWT.AccessSecret == "" {
		return fmt.Errorf("config: JWT_ACCESS_SECRET is required")
	}
	if c.JWT.RefreshSecret == "" {
		return fmt.Errorf("config: JWT_REFRESH_SECRET is required")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}
