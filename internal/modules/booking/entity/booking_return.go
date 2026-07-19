package entity

import "time"

// BookingReturn mirrors the booking_returns table.
type BookingReturn struct {
	ID                string     `db:"id"`
	TenantID          string     `db:"tenant_id"`
	BookingID         string     `db:"booking_id"`
	BookingItemID     string     `db:"booking_item_id"`
	ReturnedAt        time.Time  `db:"returned_at"`
	ConditionOnReturn string     `db:"condition_on_return"`
	LateFee           float64    `db:"late_fee"`
	DamageFee         float64    `db:"damage_fee"`
	Notes             *string    `db:"notes"`
	Status            string     `db:"status"`
	CreatedBy         *string    `db:"created_by"`
	UpdatedBy         *string    `db:"updated_by"`
	DeletedBy         *string    `db:"deleted_by"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
	Version           int        `db:"version"`
}
