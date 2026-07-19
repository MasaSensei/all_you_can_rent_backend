// Package database owns the Postgres connection pool and defines the
// Querier interface every repository in every module is written against.
// Repositories never depend on *sqlx.DB directly — they depend on Querier,
// which both *sqlx.DB and *sqlx.Tx satisfy. That is what lets a repository
// be called standalone or inside pkg/transaction.WithTx with no code
// change and no duplicated repository implementations.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Config holds Postgres connection parameters.
type Config struct {
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

// DSN builds the Postgres connection string from Config.
func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}

// Querier is the minimal surface repositories depend on. *sqlx.DB and
// *sqlx.Tx both implement it.
type Querier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
}

// New opens the connection pool and configures pool limits. It does not
// verify connectivity beyond what sqlx.Connect performs (Open + Ping).
func New(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("database: connect: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	return db, nil
}

// Health pings the database, bounded by the caller's context deadline.
func Health(ctx context.Context, db *sqlx.DB) error {
	return db.PingContext(ctx)
}
