package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos/internal/modules/auth/entity"
	"rentos/internal/modules/auth/repository"
	"rentos/pkg/database"
)

type sessionRepository struct {
	qCreate             string
	qFindByRefreshToken string
	qRevoke             string
	qRevokeAllForUser   string
}

func NewSessionRepository(qCreate, qFindByRefreshToken, qRevoke, qRevokeAllForUser string) repository.SessionRepository {
	return &sessionRepository{
		qCreate:             qCreate,
		qFindByRefreshToken: qFindByRefreshToken,
		qRevoke:             qRevoke,
		qRevokeAllForUser:   qRevokeAllForUser,
	}
}

func (r *sessionRepository) Create(ctx context.Context, q database.Querier, s *entity.UserSession) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		s.ID, s.TenantID, s.UserID, s.Token, s.RefreshToken,
		s.IPAddress, s.UserAgent, s.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("sessionRepository.Create: %w", err)
	}
	return nil
}

func (r *sessionRepository) FindByRefreshToken(ctx context.Context, q database.Querier, refreshToken string) (*entity.UserSession, error) {
	var s entity.UserSession
	if err := q.GetContext(ctx, &s, r.qFindByRefreshToken, refreshToken); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("sessionRepository.FindByRefreshToken: %w", err)
	}
	return &s, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, q database.Querier, id string) error {
	_, err := q.ExecContext(ctx, r.qRevoke, id)
	if err != nil {
		return fmt.Errorf("sessionRepository.Revoke: %w", err)
	}
	return nil
}

func (r *sessionRepository) RevokeAllForUser(ctx context.Context, q database.Querier, userID string) error {
	_, err := q.ExecContext(ctx, r.qRevokeAllForUser, userID)
	if err != nil {
		return fmt.Errorf("sessionRepository.RevokeAllForUser: %w", err)
	}
	return nil
}
