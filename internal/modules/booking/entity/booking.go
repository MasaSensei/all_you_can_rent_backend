package entity

import "time"

// Booking mirrors the bookings table.
type Booking struct {
	ID            string     `db:"id"`
	TenantID      string     `db:"tenant_id"`
	CustomerID    string     `db:"customer_id"`
	CouponID      *string    `db:"coupon_id"`
	BookingNumber string     `db:"booking_number"`
	StartDate     time.Time  `db:"start_date"`
	EndDate       time.Time  `db:"end_date"`
	Subtotal      float64    `db:"subtotal"`
	TaxTotal      float64    `db:"tax_total"`
	DiscountTotal float64    `db:"discount_total"`
	TotalAmount   float64    `db:"total_amount"`
	BookingStatus string     `db:"booking_status"`
	PaymentStatus string     `db:"payment_status"`
	Notes         *string    `db:"notes"`
	Status        string     `db:"status"`
	CreatedBy     *string    `db:"created_by"`
	UpdatedBy     *string    `db:"updated_by"`
	DeletedBy     *string    `db:"deleted_by"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
	Version       int        `db:"version"`
}

// Valid booking status transitions:
//   pending → confirmed → active → completed
//   pending → cancelled
//   confirmed → cancelled
const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusActive    = "active"
	BookingStatusCompleted = "completed"
	BookingStatusCancelled = "cancelled"
)
