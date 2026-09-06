package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos-backend/internal/modules/auth/entity"
	"rentos-backend/internal/modules/auth/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// userRepository
// ============================================================

type userRepository struct {
	qCreate         string
	qFindByID       string
	qFindByEmail    string
	qFindByUsername string
	qUpdatePassword string
}

func NewUserRepository(qCreate, qFindByID, qFindByEmail, qFindByUsername, qUpdatePassword string) repository.UserRepository {
	return &userRepository{
		qCreate: qCreate, qFindByID: qFindByID, qFindByEmail: qFindByEmail,
		qFindByUsername: qFindByUsername, qUpdatePassword: qUpdatePassword,
	}
}

func (r *userRepository) Create(ctx context.Context, q database.Querier, u *entity.User) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		u.ID, u.TenantID, u.Username, u.Email, u.PasswordHash,
		u.FirstName, u.LastName, u.IsActive, u.Status, u.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("userRepository.Create: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, q database.Querier, id string) (*entity.User, error) {
	var u entity.User
	if err := q.GetContext(ctx, &u, r.qFindByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByID: %w", err)
	}
	return &u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, q database.Querier, email string) (*entity.User, error) {
	var u entity.User
	if err := q.GetContext(ctx, &u, r.qFindByEmail, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByEmail: %w", err)
	}
	return &u, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, q database.Querier, username, tenantID string) (*entity.User, error) {
	var u entity.User
	if err := q.GetContext(ctx, &u, r.qFindByUsername, username, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByUsername: %w", err)
	}
	return &u, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, q database.Querier, id, hash string) error {
	_, err := q.ExecContext(ctx, r.qUpdatePassword, id, hash)
	return err
}

// ============================================================
// sessionRepository
// ============================================================

type sessionRepository struct {
	qCreate      string
	qFindByToken string
	qRevoke      string
}

func NewSessionRepository(qCreate, qFindByToken, qRevoke string) repository.SessionRepository {
	return &sessionRepository{qCreate: qCreate, qFindByToken: qFindByToken, qRevoke: qRevoke}
}

func (r *sessionRepository) Create(ctx context.Context, q database.Querier, s *entity.UserSession) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		s.ID, s.TenantID, s.UserID, s.RefreshToken, s.UserAgent, s.IPAddress, s.ExpiresAt,
	)
	return err
}

func (r *sessionRepository) FindByRefreshToken(ctx context.Context, q database.Querier, token string) (*entity.UserSession, error) {
	var s entity.UserSession
	if err := q.GetContext(ctx, &s, r.qFindByToken, token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("sessionRepository.FindByRefreshToken: %w", err)
	}
	return &s, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, q database.Querier, token string) error {
	_, err := q.ExecContext(ctx, r.qRevoke, token)
	return err
}

// ============================================================
// passwordResetRepository
// ============================================================

type passwordResetRepository struct {
	qCreate     string
	qFindByHash string
	qMarkUsed   string
}

func NewPasswordResetRepository(qCreate, qFindByHash, qMarkUsed string) repository.PasswordResetRepository {
	return &passwordResetRepository{qCreate: qCreate, qFindByHash: qFindByHash, qMarkUsed: qMarkUsed}
}

func (r *passwordResetRepository) Create(ctx context.Context, q database.Querier, pr *entity.PasswordReset) error {
	_, err := q.ExecContext(ctx, r.qCreate, pr.ID, pr.TenantID, pr.UserID, pr.TokenHash, pr.ExpiresAt)
	return err
}

func (r *passwordResetRepository) FindByTokenHash(ctx context.Context, q database.Querier, hash string) (*entity.PasswordReset, error) {
	var pr entity.PasswordReset
	if err := q.GetContext(ctx, &pr, r.qFindByHash, hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("passwordResetRepository.FindByTokenHash: %w", err)
	}
	return &pr, nil
}

func (r *passwordResetRepository) MarkUsed(ctx context.Context, q database.Querier, id string) error {
	_, err := q.ExecContext(ctx, r.qMarkUsed, id)
	return err
}
