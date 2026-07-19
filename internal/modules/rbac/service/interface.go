package service

import (
	"context"

	"rentos/internal/modules/rbac/dto/request"
	"rentos/internal/modules/rbac/entity"
)

// RoleService manages roles and their permission assignments.
type RoleService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateRole) (*entity.Role, error)
	GetByID(ctx context.Context, id, tenantID string) (*entity.Role, error)
	List(ctx context.Context, tenantID string) ([]entity.Role, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateRole) (*entity.Role, error)
	Delete(ctx context.Context, id, tenantID string) error
	SyncPermissions(ctx context.Context, roleID, tenantID string, req request.SyncPermissions) ([]entity.Permission, error)
	GetRolePermissions(ctx context.Context, roleID, tenantID string) ([]entity.Permission, error)
}

// PermissionService reads the global permission catalogue.
type PermissionService interface {
	List(ctx context.Context) ([]entity.Permission, error)
}

// UserRoleService manages user ↔ role assignments and permission checks.
type UserRoleService interface {
	Assign(ctx context.Context, userID, tenantID string, req request.AssignRole) error
	Revoke(ctx context.Context, userID, roleID, tenantID string) error
	ListForUser(ctx context.Context, userID, tenantID string) ([]entity.Role, error)

	// GetUserPermissionNames returns a flat list of permission name strings
	// (e.g. "booking.create") for the given user. Used by RBACMiddleware.
	GetUserPermissionNames(ctx context.Context, userID, tenantID string) ([]string, error)
}
