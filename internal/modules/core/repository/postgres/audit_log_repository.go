package postgres

import (
	"context"
	"fmt"

	"rentos/internal/modules/core/entity"
	"rentos/internal/modules/core/repository"
	"rentos/pkg/database"
)

type auditLogRepository struct {
	qInsert string
	qList   string
}

// NewAuditLogRepository builds the Postgres-backed AuditLogRepository
// from pre-loaded query strings.
func NewAuditLogRepository(qInsert, qList string) repository.AuditLogRepository {
	return &auditLogRepository{qInsert: qInsert, qList: qList}
}

func (r *auditLogRepository) Create(ctx context.Context, q database.Querier, a *entity.AuditLog) error {
	_, err := q.ExecContext(ctx, r.qInsert,
		a.ID, a.TenantID, a.UserID, a.EntityType, a.EntityID, a.Action,
		a.OldValues, a.NewValues, a.IPAddress, a.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("auditLogRepository.Create: %w", err)
	}
	return nil
}

func (r *auditLogRepository) List(ctx context.Context, q database.Querier, tenantID string, entityType, entityID, action *string, limit, offset int) ([]entity.AuditLog, error) {
	var out []entity.AuditLog
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, entityType, entityID, action, limit, offset); err != nil {
		return nil, fmt.Errorf("auditLogRepository.List: %w", err)
	}
	return out, nil
}
