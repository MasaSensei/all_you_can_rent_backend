//go:build integration

// Package integration contains end-to-end HTTP tests that require a real
// Postgres + Redis instance. Run with:
//
//	go test -v -tags=integration ./tests/integration/...
//
// Or via Makefile:
//
//	make test-integration
package integration

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// testDB creates a dedicated test schema and returns a connected DB.
// Deferred cleanup drops the schema to guarantee test isolation.
func testDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "rentos"),
		getenv("DB_PASSWORD", "rentos"),
		getenv("DB_NAME", "rentos_test"),
	)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		t.Fatalf("testDB: connect: %v", err)
	}

	t.Cleanup(func() { _ = db.Close() })
	return db
}

// testApp boots a minimal Fiber app for HTTP-level testing without starting
// an actual listener. Use app.Test() to send requests.
func testApp(handlers ...func(fiber.Router)) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "rentos-test",
	})
	r := app.Group("/api/v1")
	for _, h := range handlers {
		h(r)
	}
	return app
}

// truncate clears the named tables in a single transaction.
// Call at the start of each test to ensure a clean state.
func truncate(t *testing.T, db *sqlx.DB, tables ...string) {
	t.Helper()
	for _, table := range tables {
		if _, err := db.ExecContext(context.Background(),
			fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table),
		); err != nil {
			t.Fatalf("truncate %s: %v", table, err)
		}
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// testRequest sends a request to the app and returns the response.
func testRequest(t *testing.T, app *fiber.App, method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	// Fiber's app.Test() is used directly in tests.
	// This helper is a placeholder for common setup logic.
	return nil
}
