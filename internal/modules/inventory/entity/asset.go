package entity

import (
	"time"
)

// Asset mirrors the assets table.
type Asset struct {
	ID               string     `db:"id"`
	TenantID         string     `db:"tenant_id"`
	CategoryID       *string    `db:"category_id"`
	AssetTemplateID  *string    `db:"asset_template_id"`
	Name             string     `db:"name"`
	SKU              *string    `db:"sku"`
	SerialNumber     *string    `db:"serial_number"`
	Description      *string    `db:"description"`
	PurchasePrice    *float64   `db:"purchase_price"`
	ReplacementValue *float64   `db:"replacement_value"`
	PurchaseDate     *time.Time `db:"purchase_date"`
	Condition        string     `db:"condition"`
	Location         *string    `db:"location"`
	Status           string     `db:"status"`
	CreatedBy        *string    `db:"created_by"`
	UpdatedBy        *string    `db:"updated_by"`
	DeletedBy        *string    `db:"deleted_by"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
	Version          int        `db:"version"`
}
