package request

// CreateRole is the input for creating a new role.
type CreateRole struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description *string `json:"description" validate:"omitempty,max=255"`
}

// UpdateRole is the input for updating an existing role.
type UpdateRole struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description" validate:"omitempty,max=255"`
}

// SyncPermissions replaces a role's entire permission set.
type SyncPermissions struct {
	PermissionIDs []string `json:"permission_ids" validate:"required,dive,uuid"`
}

// AssignRole assigns a role to a user.
type AssignRole struct {
	RoleID string `json:"role_id" validate:"required,uuid"`
}
