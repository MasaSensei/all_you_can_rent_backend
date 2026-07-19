package response

import "time"

// Permission is the API-facing shape of entity.Permission.
type Permission struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Module string `json:"module"`
	Action string `json:"action"`
}

// Role is the API-facing shape of entity.Role.
type Role struct {
	ID          string       `json:"id"`
	TenantID    string       `json:"tenant_id"`
	Name        string       `json:"name"`
	Description *string      `json:"description,omitempty"`
	IsSystem    bool         `json:"is_system"`
	Status      string       `json:"status"`
	Permissions []Permission `json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
