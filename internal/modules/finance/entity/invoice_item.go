package entity

import "time"

// InvoiceItem mirrors the invoice_items table.
type InvoiceItem struct {
	ID            string     `db:"id"`
	TenantID      string     `db:"tenant_id"`
	InvoiceID     string     `db:"invoice_id"`
	BookingItemID *string    `db:"booking_item_id"`
	Description   string     `db:"description"`
	Quantity      int        `db:"quantity"`
	UnitPrice     float64    `db:"unit_price"`
	TaxAmount     float64    `db:"tax_amount"`
	LineTotal     float64    `db:"line_total"`
	Status        string     `db:"status"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
	Version       int        `db:"version"`
}
