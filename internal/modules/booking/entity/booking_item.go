package entity

import "time"

// BookingItem mirrors the booking_items table.
type BookingItem struct {
	ID        string     `db:"id"`
	TenantID  string     `db:"tenant_id"`
	BookingID string     `db:"booking_id"`
	AssetID   string     `db:"asset_id"`
	Quantity  int        `db:"quantity"`
	UnitPrice float64    `db:"unit_price"`
	LineTotal float64    `db:"line_total"`
	StartDate time.Time  `db:"start_date"`
	EndDate   time.Time  `db:"end_date"`
	Status    string     `db:"status"`
	CreatedBy *string    `db:"created_by"`
	UpdatedBy *string    `db:"updated_by"`
	DeletedBy *string    `db:"deleted_by"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Version   int        `db:"version"`
}
