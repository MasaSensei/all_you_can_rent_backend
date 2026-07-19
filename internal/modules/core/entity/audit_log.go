package entity

import "time"

// AuditLog mirrors the audit_logs table.
type AuditLog struct {
	ID         string     `db:"id"`
	TenantID   string     `db:"tenant_id"`
	UserID     *string    `db:"user_id"`
	EntityType string     `db:"entity_type"`
	EntityID   string     `db:"entity_id"`
	Action     string     `db:"action"`
	OldValues  *string    `db:"old_values"`
	NewValues  *string    `db:"new_values"`
	IPAddress  *string    `db:"ip_address"`
	UserAgent  *string    `db:"user_agent"`
	Status     string     `db:"status"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
	Version    int        `db:"version"`
}
