package repository

import (
	"context"
	"errors"

	"rentos/internal/modules/rbac/entity"
	"rentos/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

// RoleRepository manages roles and their permission sets.
type RoleRepository interface {
	Create(ctx context.Context, q database.Querier, r *entity.Role) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Role, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Role, error)
	Update(ctx context.Context, q database.Querier, r *entity.Role) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// PermissionRepository reads the global permissions catalogue.
type PermissionRepository interface {
	List(ctx context.Context, q database.Querier) ([]entity.Permission, error)
	FindByIDs(ctx context.Context, q database.Querier, ids []string) ([]entity.Permission, error)
	GetByRole(ctx context.Context, q database.Querier, roleID, tenantID string) ([]entity.Permission, error)
}

// RolePermissionRepository manages the many-to-many between roles and permissions.
type RolePermissionRepository interface {
	Sync(ctx context.Context, q database.Querier, roleID, tenantID string, permissionIDs []string) error
}

// UserRoleRepository manages the many-to-many between users and roles.
type UserRoleRepository interface {
	Assign(ctx context.Context, q database.Querier, ur *entity.UserRole) error
	Revoke(ctx context.Context, q database.Querier, userID, roleID, tenantID string) error
	ListForUser(ctx context.Context, q database.Querier, userID, tenantID string) ([]entity.Role, error)
	GetUserPermissionNames(ctx context.Context, q database.Querier, userID, tenantID string) ([]string, error)
}
