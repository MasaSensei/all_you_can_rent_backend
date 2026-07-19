// Package repository defines the contracts the auth service depends on.
package repository

import (
	"context"
	"errors"

	"rentos/internal/modules/auth/entity"
	"rentos/pkg/database"
)

// ErrNotFound is returned by any auth repository when a queried row
// does not exist.
var ErrNotFound = errors.New("repository: record not found")

// UserRepository handles persistence for the users table.
type UserRepository interface {
	Create(ctx context.Context, q database.Querier, u *entity.User) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.User, error)
	FindByEmail(ctx context.Context, q database.Querier, email, tenantID string) (*entity.User, error)
	List(ctx context.Context, q database.Querier, tenantID string, limit, offset int) ([]entity.User, error)
	Update(ctx context.Context, q database.Querier, u *entity.User) error
	UpdateLastLogin(ctx context.Context, q database.Querier, id string) error
	UpdatePassword(ctx context.Context, q database.Querier, id, hash string) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// SessionRepository handles persistence for the user_sessions table.
type SessionRepository interface {
	Create(ctx context.Context, q database.Querier, s *entity.UserSession) error
	FindByRefreshToken(ctx context.Context, q database.Querier, refreshToken string) (*entity.UserSession, error)
	Revoke(ctx context.Context, q database.Querier, id string) error
	RevokeAllForUser(ctx context.Context, q database.Querier, userID string) error
}

// PasswordResetRepository handles persistence for the password_resets table.
type PasswordResetRepository interface {
	Create(ctx context.Context, q database.Querier, pr *entity.PasswordReset) error
	FindByToken(ctx context.Context, q database.Querier, token string) (*entity.PasswordReset, error)
	Consume(ctx context.Context, q database.Querier, id string) error
}
