package entity

import "time"

// AssetTemplate mirrors the asset_templates table.
type AssetTemplate struct {
	ID          string     `db:"id"`
	TenantID    string     `db:"tenant_id"`
	CategoryID  *string    `db:"category_id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	Status      string     `db:"status"`
	CreatedBy   *string    `db:"created_by"`
	UpdatedBy   *string    `db:"updated_by"`
	DeletedBy   *string    `db:"deleted_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Version     int        `db:"version"`
}
