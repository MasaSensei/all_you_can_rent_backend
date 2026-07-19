package entity

import "time"

// Setting mirrors the settings table (per-tenant key/value configuration).
type Setting struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	Key       string     `db:"key"`
	Value     string     `db:"value"`
	Type      string     `db:"type"`
	Status    string     `db:"status"`
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}
