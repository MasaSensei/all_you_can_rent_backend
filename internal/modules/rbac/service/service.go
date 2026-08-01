package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/rbac/dto/request"
	"rentos-backend/internal/modules/rbac/entity"
	"rentos-backend/internal/modules/rbac/repository"
	"rentos-backend/pkg/response"
	"rentos-backend/pkg/transaction"
)

// ============================================================
// roleService
// ============================================================

type roleService struct {
	db          *sqlx.DB
	roles       repository.RoleRepository
	permissions repository.PermissionRepository
	rolePerms   repository.RolePermissionRepository
}

func NewRoleService(
	db *sqlx.DB,
	roles repository.RoleRepository,
	permissions repository.PermissionRepository,
	rolePerms repository.RolePermissionRepository,
) RoleService {
	return &roleService{db: db, roles: roles, permissions: permissions, rolePerms: rolePerms}
}

func (s *roleService) Create(ctx context.Context, tenantID, actorID string, req request.CreateRole) (*entity.Role, error) {
	r := &entity.Role{
		ID:          uuid.NewString(),
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   &actorID,
	}
	if err := s.roles.Create(ctx, s.db, r); err != nil {
		return nil, err
	}
	return s.roles.FindByID(ctx, s.db, r.ID, tenantID)
}

func (s *roleService) GetByID(ctx context.Context, id, tenantID string) (*entity.Role, error) {
	r, err := s.roles.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.NewAppError(response.CodeNotFound, "role not found")
		}
		return nil, err
	}
	return r, nil
}

func (s *roleService) List(ctx context.Context, tenantID string) ([]entity.Role, error) {
	return s.roles.List(ctx, s.db, tenantID)
}

func (s *roleService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateRole) (*entity.Role, error) {
	r := &entity.Role{
		ID:          id,
		TenantID:    tenantID,
		Name:        *req.Name,
		Description: req.Description,
		UpdatedBy:   &actorID,
	}
	if err := s.roles.Update(ctx, s.db, r); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.NewAppError(response.CodeNotFound, "role not found or is a system role")
		}
		return nil, err
	}
	return s.roles.FindByID(ctx, s.db, id, tenantID)
}

func (s *roleService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.roles.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return response.NewAppError(response.CodeNotFound, "role not found or is a system role")
		}
		return err
	}
	return nil
}

// SyncPermissions replaces the role's permission set atomically.
func (s *roleService) SyncPermissions(ctx context.Context, roleID, tenantID string, req request.SyncPermissions) ([]entity.Permission, error) {
	// Validate role exists.
	if _, err := s.GetByID(ctx, roleID, tenantID); err != nil {
		return nil, err
	}

	// Validate all supplied permission IDs exist.
	if len(req.PermissionIDs) > 0 {
		found, err := s.permissions.FindByIDs(ctx, s.db, req.PermissionIDs)
		if err != nil {
			return nil, err
		}
		if len(found) != len(req.PermissionIDs) {
			return nil, response.NewAppError(response.CodeValidation, "one or more permission IDs are invalid")
		}
	}

	if err := transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		return s.rolePerms.Sync(ctx, tx, roleID, tenantID, req.PermissionIDs)
	}); err != nil {
		return nil, err
	}

	return s.permissions.GetByRole(ctx, s.db, roleID, tenantID)
}

func (s *roleService) GetRolePermissions(ctx context.Context, roleID, tenantID string) ([]entity.Permission, error) {
	if _, err := s.GetByID(ctx, roleID, tenantID); err != nil {
		return nil, err
	}
	return s.permissions.GetByRole(ctx, s.db, roleID, tenantID)
}

// ============================================================
// permissionService
// ============================================================

type permissionService struct {
	db   *sqlx.DB
	repo repository.PermissionRepository
}

func NewPermissionService(db *sqlx.DB, repo repository.PermissionRepository) PermissionService {
	return &permissionService{db: db, repo: repo}
}

func (s *permissionService) List(ctx context.Context) ([]entity.Permission, error) {
	return s.repo.List(ctx, s.db)
}

// ============================================================
// userRoleService
// ============================================================

type userRoleService struct {
	db        *sqlx.DB
	userRoles repository.UserRoleRepository
	roles     repository.RoleRepository
}

func NewUserRoleService(
	db *sqlx.DB,
	userRoles repository.UserRoleRepository,
	roles repository.RoleRepository,
) UserRoleService {
	return &userRoleService{db: db, userRoles: userRoles, roles: roles}
}

func (s *userRoleService) Assign(ctx context.Context, userID, tenantID string, req request.AssignRole) error {
	// Validate role belongs to tenant.
	if _, err := s.roles.FindByID(ctx, s.db, req.RoleID, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return response.NewAppError(response.CodeNotFound, "role not found")
		}
		return err
	}

	ur := &entity.UserRole{
		ID:       uuid.NewString(),
		TenantID: tenantID,
		UserID:   userID,
		RoleID:   req.RoleID,
	}
	return s.userRoles.Assign(ctx, s.db, ur)
}

func (s *userRoleService) Revoke(ctx context.Context, userID, roleID, tenantID string) error {
	if err := s.userRoles.Revoke(ctx, s.db, userID, roleID, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return response.NewAppError(response.CodeNotFound, "user role assignment not found")
		}
		return err
	}
	return nil
}

func (s *userRoleService) ListForUser(ctx context.Context, userID, tenantID string) ([]entity.Role, error) {
	return s.userRoles.ListForUser(ctx, s.db, userID, tenantID)
}

func (s *userRoleService) GetUserPermissionNames(ctx context.Context, userID, tenantID string) ([]string, error) {
	return s.userRoles.GetUserPermissionNames(ctx, s.db, userID, tenantID)
}
