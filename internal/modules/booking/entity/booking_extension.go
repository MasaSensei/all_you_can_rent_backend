package entity

import "time"

// BookingExtension mirrors the booking_extensions table.
type BookingExtension struct {
	ID             string     `db:"id"`
	TenantID       string     `db:"tenant_id"`
	BookingID      string     `db:"booking_id"`
	BookingItemID  string     `db:"booking_item_id"`
	OldEndDate     time.Time  `db:"old_end_date"`
	NewEndDate     time.Time  `db:"new_end_date"`
	AdditionalCost float64    `db:"additional_cost"`
	Reason         *string    `db:"reason"`
	Status         string     `db:"status"`
	CreatedBy      *string    `db:"created_by"`
	UpdatedBy      *string    `db:"updated_by"`
	DeletedBy      *string    `db:"deleted_by"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
	Version        int        `db:"version"`
}
