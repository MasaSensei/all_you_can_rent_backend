package repository

import (
	"context"
	"errors"

	"rentos-backend/internal/modules/auth/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

type UserRepository interface {
	Create(ctx context.Context, q database.Querier, u *entity.User) error
	FindByID(ctx context.Context, q database.Querier, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, q database.Querier, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, q database.Querier, username, tenantID string) (*entity.User, error)
	UpdatePassword(ctx context.Context, q database.Querier, id, hash string) error
}

type SessionRepository interface {
	Create(ctx context.Context, q database.Querier, s *entity.UserSession) error
	FindByRefreshToken(ctx context.Context, q database.Querier, token string) (*entity.UserSession, error)
	Revoke(ctx context.Context, q database.Querier, token string) error
}

type PasswordResetRepository interface {
	Create(ctx context.Context, q database.Querier, pr *entity.PasswordReset) error
	FindByTokenHash(ctx context.Context, q database.Querier, hash string) (*entity.PasswordReset, error)
	MarkUsed(ctx context.Context, q database.Querier, id string) error
}
