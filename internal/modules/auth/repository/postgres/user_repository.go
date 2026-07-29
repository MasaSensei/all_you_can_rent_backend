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

type userRepository struct {
	qCreate          string
	qFindByID        string
	qFindByEmail     string
	qList            string
	qUpdate          string
	qUpdateLastLogin string
	qUpdatePassword  string
	qDelete          string
}

func NewUserRepository(
	qCreate, qFindByID, qFindByEmail, qList,
	qUpdate, qUpdateLastLogin, qUpdatePassword, qDelete string,
) repository.UserRepository {
	return &userRepository{
		qCreate:          qCreate,
		qFindByID:        qFindByID,
		qFindByEmail:     qFindByEmail,
		qList:            qList,
		qUpdate:          qUpdate,
		qUpdateLastLogin: qUpdateLastLogin,
		qUpdatePassword:  qUpdatePassword,
		qDelete:          qDelete,
	}
}

func (r *userRepository) Create(ctx context.Context, q database.Querier, u *entity.User) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		u.ID, u.TenantID, u.Email, u.PasswordHash,
		u.FirstName, u.LastName, u.Phone,
	)
	if err != nil {
		return fmt.Errorf("userRepository.Create: %w", err)
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.User, error) {
	var u entity.User
	if err := q.GetContext(ctx, &u, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByID: %w", err)
	}
	return &u, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, q database.Querier, email, tenantID string) (*entity.User, error) {
	var u entity.User
	if err := q.GetContext(ctx, &u, r.qFindByEmail, email, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("userRepository.FindByEmail: %w", err)
	}
	return &u, nil
}

func (r *userRepository) List(ctx context.Context, q database.Querier, tenantID string, limit, offset int) ([]entity.User, error) {
	var out []entity.User
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, limit, offset); err != nil {
		return nil, fmt.Errorf("userRepository.List: %w", err)
	}
	return out, nil
}

func (r *userRepository) Update(ctx context.Context, q database.Querier, u *entity.User) error {
	res, err := q.ExecContext(ctx, r.qUpdate,
		u.ID, u.TenantID, u.FirstName, u.LastName, u.Phone,
	)
	if err != nil {
		return fmt.Errorf("userRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, q database.Querier, id string) error {
	_, err := q.ExecContext(ctx, r.qUpdateLastLogin, id)
	if err != nil {
		return fmt.Errorf("userRepository.UpdateLastLogin: %w", err)
	}
	return nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, q database.Querier, id, hash string) error {
	res, err := q.ExecContext(ctx, r.qUpdatePassword, id, hash)
	if err != nil {
		return fmt.Errorf("userRepository.UpdatePassword: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("userRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}
