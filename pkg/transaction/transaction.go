// Package transaction provides the single shared helper services use to
// run multiple repository writes atomically. Services depend on this
// instead of each reimplementing begin/commit/rollback handling.
package transaction

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Fn is a unit of work executed inside one database transaction. It
// receives a *sqlx.Tx, which satisfies database.Querier, so repository
// methods accept it exactly as they would accept the plain *sqlx.DB.
type Fn func(tx *sqlx.Tx) error

// WithTx begins a transaction, runs fn, and commits or rolls back based on
// the returned error. A panic inside fn is rolled back and re-panicked so
// it still reaches the global recover middleware.
func WithTx(ctx context.Context, db *sqlx.DB, fn Fn) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("transaction: begin: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction: rollback failed: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction: commit: %w", err)
	}
	return nil
}
