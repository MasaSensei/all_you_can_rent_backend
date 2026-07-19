package entity

import "time"

// Category mirrors the categories table.
type Category struct {
	ID          string     `db:"id"`
	TenantID    string     `db:"tenant_id"`
	ParentID    *string    `db:"parent_id"`
	Name        string     `db:"name"`
	Slug        string     `db:"slug"`
	Description *string    `db:"description"`
	Icon        *string    `db:"icon"`
	SortOrder   int        `db:"sort_order"`
	Status      string     `db:"status"`
	CreatedBy   *string    `db:"created_by"`
	UpdatedBy   *string    `db:"updated_by"`
	DeletedBy   *string    `db:"deleted_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Version     int        `db:"version"`
}
