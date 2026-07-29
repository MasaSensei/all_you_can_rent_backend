package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos-backend/internal/modules/core/entity"
	"rentos-backend/internal/modules/core/repository"
	"rentos-backend/pkg/database"
)

type settingRepository struct {
	qUpsert string
	qFind   string
	qList   string
}

// NewSettingRepository builds the Postgres-backed SettingRepository from
// pre-loaded query strings.
func NewSettingRepository(qUpsert, qFind, qList string) repository.SettingRepository {
	return &settingRepository{qUpsert: qUpsert, qFind: qFind, qList: qList}
}

func (r *settingRepository) Upsert(ctx context.Context, q database.Querier, s *entity.Setting) error {
	_, err := q.ExecContext(ctx, r.qUpsert, s.ID, s.TenantID, s.Key, s.Value, s.Type, s.UpdatedBy)
	if err != nil {
		return fmt.Errorf("settingRepository.Upsert: %w", err)
	}
	return nil
}

func (r *settingRepository) Find(ctx context.Context, q database.Querier, tenantID, key string) (*entity.Setting, error) {
	var s entity.Setting
	if err := q.GetContext(ctx, &s, r.qFind, tenantID, key); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("settingRepository.Find: %w", err)
	}
	return &s, nil
}

func (r *settingRepository) ListByTenant(ctx context.Context, q database.Querier, tenantID string) ([]entity.Setting, error) {
	var out []entity.Setting
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("settingRepository.ListByTenant: %w", err)
	}
	return out, nil
}
