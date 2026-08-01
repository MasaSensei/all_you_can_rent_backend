package entity

import "time"

// Tax mirrors the taxes table.
type Tax struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	Name      string     `db:"name"`
	Rate      float64    `db:"rate"`
	TaxType   string     `db:"tax_type"` // percentage, fixed
	IsDefault bool       `db:"is_default"`
	Status    string     `db:"status"`
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}

const (
	TaxTypePercentage = "percentage"
	TaxTypeFixed      = "fixed"
)
