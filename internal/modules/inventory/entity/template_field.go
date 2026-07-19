package entity

import "time"

// TemplateField mirrors the template_fields table.
type TemplateField struct {
	ID              string     `db:"id"`
	TenantID        string     `db:"tenant_id"`
	AssetTemplateID string     `db:"asset_template_id"`
	FieldName       string     `db:"field_name"`
	FieldLabel      string     `db:"field_label"`
	FieldType       string     `db:"field_type"`
	IsRequired      bool       `db:"is_required"`
	DefaultValue    *string    `db:"default_value"`
	Options         *string    `db:"options"` // JSON stored as string
	SortOrder       int        `db:"sort_order"`
	Status          string     `db:"status"`
	CreatedBy       *string    `db:"created_by"`
	UpdatedBy       *string    `db:"updated_by"`
	DeletedBy       *string    `db:"deleted_by"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
	Version         int        `db:"version"`
}
