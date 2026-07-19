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

type tenantRepository struct {
	qCreate       string
	qFindByID     string
	qFindBySlug   string
	qUpdateStatus string
}

// NewTenantRepository builds the Postgres-backed TenantRepository from
// pre-loaded query strings.
func NewTenantRepository(qCreate, qFindByID, qFindBySlug, qUpdateStatus string) repository.TenantRepository {
	return &tenantRepository{
		qCreate:       qCreate,
		qFindByID:     qFindByID,
		qFindBySlug:   qFindBySlug,
		qUpdateStatus: qUpdateStatus,
	}
}

func (r *tenantRepository) Create(ctx context.Context, q database.Querier, t *entity.Tenant) error {
	_, err := q.ExecContext(ctx, r.qCreate, t.ID, t.Name, t.Slug, t.Domain, t.Plan, t.Status)
	if err != nil {
		return fmt.Errorf("tenantRepository.Create: %w", err)
	}
	return nil
}

func (r *tenantRepository) FindByID(ctx context.Context, q database.Querier, id string) (*entity.Tenant, error) {
	var t entity.Tenant
	if err := q.GetContext(ctx, &t, r.qFindByID, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("tenantRepository.FindByID: %w", err)
	}
	return &t, nil
}

func (r *tenantRepository) FindBySlug(ctx context.Context, q database.Querier, slug string) (*entity.Tenant, error) {
	var t entity.Tenant
	if err := q.GetContext(ctx, &t, r.qFindBySlug, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("tenantRepository.FindBySlug: %w", err)
	}
	return &t, nil
}

func (r *tenantRepository) UpdateStatus(ctx context.Context, q database.Querier, id, status string) error {
	res, err := q.ExecContext(ctx, r.qUpdateStatus, id, status)
	if err != nil {
		return fmt.Errorf("tenantRepository.UpdateStatus: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}
