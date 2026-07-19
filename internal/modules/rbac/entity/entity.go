package entity

import "time"

// Role mirrors the roles table.
type Role struct {
	ID          string     `db:"id"`
	TenantID    string     `db:"tenant_id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	IsSystem    bool       `db:"is_system"`
	Status      string     `db:"status"`
	CreatedBy   *string    `db:"created_by"`
	UpdatedBy   *string    `db:"updated_by"`
	DeletedBy   *string    `db:"deleted_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Version     int        `db:"version"`
}

// Permission mirrors the permissions table.
type Permission struct {
	ID          string     `db:"id"`
	Name        string     `db:"name"`
	Module      string     `db:"module"`
	Action      string     `db:"action"`
	Description *string    `db:"description"`
	Status      string     `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Version     int        `db:"version"`
}

// RolePermission mirrors the role_permissions table.
type RolePermission struct {
	ID           string     `db:"id"`
	TenantID     string     `db:"tenant_id"`
	RoleID       string     `db:"role_id"`
	PermissionID string     `db:"permission_id"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	Version      int        `db:"version"`
}

// UserRole mirrors the user_roles table.
type UserRole struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	UserID    string     `db:"user_id"`
	RoleID    string     `db:"role_id"`
	Status    string     `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}
