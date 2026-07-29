package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/rbac/entity"
	"rentos-backend/internal/modules/rbac/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// roleRepository
// ============================================================

type roleRepository struct {
	qCreate   string
	qFindByID string
	qList     string
	qUpdate   string
	qDelete   string
}

func NewRoleRepository(qCreate, qFindByID, qList, qUpdate, qDelete string) repository.RoleRepository {
	return &roleRepository{qCreate: qCreate, qFindByID: qFindByID, qList: qList, qUpdate: qUpdate, qDelete: qDelete}
}

func (r *roleRepository) Create(ctx context.Context, q database.Querier, role *entity.Role) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		role.ID, role.TenantID, role.Name, role.Description, role.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("roleRepository.Create: %w", err)
	}
	return nil
}

func (r *roleRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Role, error) {
	var role entity.Role
	if err := q.GetContext(ctx, &role, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("roleRepository.FindByID: %w", err)
	}
	return &role, nil
}

func (r *roleRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Role, error) {
	var out []entity.Role
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("roleRepository.List: %w", err)
	}
	return out, nil
}

func (r *roleRepository) Update(ctx context.Context, q database.Querier, role *entity.Role) error {
	res, err := q.ExecContext(ctx, r.qUpdate,
		role.ID, role.TenantID, role.Name, role.Description,
	)
	if err != nil {
		return fmt.Errorf("roleRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *roleRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("roleRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// permissionRepository
// ============================================================

type permissionRepository struct {
	qList      string
	qFindByIDs string
	qGetByRole string
}

func NewPermissionRepository(qList, qFindByIDs, qGetByRole string) repository.PermissionRepository {
	return &permissionRepository{qList: qList, qFindByIDs: qFindByIDs, qGetByRole: qGetByRole}
}

func (r *permissionRepository) List(ctx context.Context, q database.Querier) ([]entity.Permission, error) {
	var out []entity.Permission
	if err := q.SelectContext(ctx, &out, r.qList); err != nil {
		return nil, fmt.Errorf("permissionRepository.List: %w", err)
	}
	return out, nil
}

func (r *permissionRepository) FindByIDs(ctx context.Context, q database.Querier, ids []string) ([]entity.Permission, error) {
	var out []entity.Permission
	if err := q.SelectContext(ctx, &out, r.qFindByIDs, ids); err != nil {
		return nil, fmt.Errorf("permissionRepository.FindByIDs: %w", err)
	}
	return out, nil
}

func (r *permissionRepository) GetByRole(ctx context.Context, q database.Querier, roleID, tenantID string) ([]entity.Permission, error) {
	var out []entity.Permission
	if err := q.SelectContext(ctx, &out, r.qGetByRole, roleID, tenantID); err != nil {
		return nil, fmt.Errorf("permissionRepository.GetByRole: %w", err)
	}
	return out, nil
}

// ============================================================
// rolePermissionRepository
// ============================================================

type rolePermissionRepository struct {
	qDeleteAll string
	qCreate    string
}

func NewRolePermissionRepository(qDeleteAll, qCreate string) repository.RolePermissionRepository {
	return &rolePermissionRepository{qDeleteAll: qDeleteAll, qCreate: qCreate}
}

// Sync replaces the role's permissions atomically using a transaction-aware
// Querier. The caller (service) is responsible for wrapping this in
// pkg/transaction.WithTx when needed.
func (r *rolePermissionRepository) Sync(ctx context.Context, q database.Querier, roleID, tenantID string, permissionIDs []string) error {
	if _, err := q.ExecContext(ctx, r.qDeleteAll, roleID, tenantID); err != nil {
		return fmt.Errorf("rolePermissionRepository.Sync(delete): %w", err)
	}
	for _, permID := range permissionIDs {
		if _, err := q.ExecContext(ctx, r.qCreate, uuid.NewString(), tenantID, roleID, permID); err != nil {
			return fmt.Errorf("rolePermissionRepository.Sync(insert %s): %w", permID, err)
		}
	}
	return nil
}

// ============================================================
// userRoleRepository
// ============================================================

type userRoleRepository struct {
	qAssign             string
	qRevoke             string
	qListForUser        string
	qGetUserPermissions string
}

func NewUserRoleRepository(qAssign, qRevoke, qListForUser, qGetUserPermissions string) repository.UserRoleRepository {
	return &userRoleRepository{
		qAssign:             qAssign,
		qRevoke:             qRevoke,
		qListForUser:        qListForUser,
		qGetUserPermissions: qGetUserPermissions,
	}
}

func (r *userRoleRepository) Assign(ctx context.Context, q database.Querier, ur *entity.UserRole) error {
	_, err := q.ExecContext(ctx, r.qAssign, ur.ID, ur.TenantID, ur.UserID, ur.RoleID)
	if err != nil {
		return fmt.Errorf("userRoleRepository.Assign: %w", err)
	}
	return nil
}

func (r *userRoleRepository) Revoke(ctx context.Context, q database.Querier, userID, roleID, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qRevoke, userID, roleID, tenantID)
	if err != nil {
		return fmt.Errorf("userRoleRepository.Revoke: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *userRoleRepository) ListForUser(ctx context.Context, q database.Querier, userID, tenantID string) ([]entity.Role, error) {
	var out []entity.Role
	if err := q.SelectContext(ctx, &out, r.qListForUser, userID, tenantID); err != nil {
		return nil, fmt.Errorf("userRoleRepository.ListForUser: %w", err)
	}
	return out, nil
}

func (r *userRoleRepository) GetUserPermissionNames(ctx context.Context, q database.Querier, userID, tenantID string) ([]string, error) {
	// sqlx SelectContext into []string requires a workaround: scan into
	// a slice of structs then pluck the field, or use sqlx.DB directly.
	// We cast q to *sqlx.DB / *sqlx.Tx using the concrete type here
	// because SelectContext on a plain []string is not supported by sqlx
	// directly — we need to scan row-by-row.
	rows, err := q.QueryxContext(ctx, r.qGetUserPermissions, userID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("userRoleRepository.GetUserPermissionNames: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("userRoleRepository.GetUserPermissionNames scan: %w", err)
		}
		names = append(names, name)
	}
	return names, rows.Err()
}

// Ensure *sqlx.DB / *sqlx.Tx satisfy database.Querier so tests can pass
// either without a type assertion.
var _ database.Querier = (*sqlx.DB)(nil)
var _ database.Querier = (*sqlx.Tx)(nil)
