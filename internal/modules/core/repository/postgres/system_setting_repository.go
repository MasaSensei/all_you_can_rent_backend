package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos/internal/modules/core/entity"
	"rentos/internal/modules/core/repository"
	"rentos/pkg/database"
)

type systemSettingRepository struct {
	qUpsert string
	qFind   string
}

// NewSystemSettingRepository builds the Postgres-backed
// SystemSettingRepository from pre-loaded query strings.
func NewSystemSettingRepository(qUpsert, qFind string) repository.SystemSettingRepository {
	return &systemSettingRepository{qUpsert: qUpsert, qFind: qFind}
}

func (r *systemSettingRepository) Upsert(ctx context.Context, q database.Querier, s *entity.SystemSetting) error {
	_, err := q.ExecContext(ctx, r.qUpsert, s.ID, s.Key, s.Value, s.Type)
	if err != nil {
		return fmt.Errorf("systemSettingRepository.Upsert: %w", err)
	}
	return nil
}

func (r *systemSettingRepository) Find(ctx context.Context, q database.Querier, key string) (*entity.SystemSetting, error) {
	var s entity.SystemSetting
	if err := q.GetContext(ctx, &s, r.qFind, key); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("systemSettingRepository.Find: %w", err)
	}
	return &s, nil
}
