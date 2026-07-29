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

type passwordResetRepository struct {
	qCreate      string
	qFindByToken string
	qConsume     string
}

func NewPasswordResetRepository(qCreate, qFindByToken, qConsume string) repository.PasswordResetRepository {
	return &passwordResetRepository{
		qCreate:      qCreate,
		qFindByToken: qFindByToken,
		qConsume:     qConsume,
	}
}

func (r *passwordResetRepository) Create(ctx context.Context, q database.Querier, pr *entity.PasswordReset) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		pr.ID, pr.TenantID, pr.UserID, pr.Token, pr.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("passwordResetRepository.Create: %w", err)
	}
	return nil
}

func (r *passwordResetRepository) FindByToken(ctx context.Context, q database.Querier, token string) (*entity.PasswordReset, error) {
	var pr entity.PasswordReset
	if err := q.GetContext(ctx, &pr, r.qFindByToken, token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("passwordResetRepository.FindByToken: %w", err)
	}
	return &pr, nil
}

func (r *passwordResetRepository) Consume(ctx context.Context, q database.Querier, id string) error {
	_, err := q.ExecContext(ctx, r.qConsume, id)
	if err != nil {
		return fmt.Errorf("passwordResetRepository.Consume: %w", err)
	}
	return nil
}
