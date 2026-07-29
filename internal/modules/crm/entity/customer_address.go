package entity

import "time"

// CustomerAddress mirrors the customer_addresses table.
type CustomerAddress struct {
	ID          string     `db:"id"`
	TenantID    string     `db:"tenant_id"`
	CustomerID  string     `db:"customer_id"`
	AddressType string     `db:"address_type"` // billing, shipping
	Line1       string     `db:"line1"`
	Line2       *string    `db:"line2"`
	City        string     `db:"city"`
	State       *string    `db:"state"`
	PostalCode  *string    `db:"postal_code"`
	Country     string     `db:"country"`
	IsDefault   bool       `db:"is_default"`
	Status      string     `db:"status"`
	CreatedBy   *string    `db:"created_by"`
	UpdatedBy   *string    `db:"updated_by"`
	DeletedBy   *string    `db:"deleted_by"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Version     int        `db:"version"`
}
