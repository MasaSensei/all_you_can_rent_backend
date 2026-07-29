package entity

import "time"

// Coupon mirrors the coupons table.
type Coupon struct {
	ID            string     `db:"id"`
	TenantID      string     `db:"tenant_id"`
	Code          string     `db:"code"`
	DiscountType  string     `db:"discount_type"`  // percentage, fixed
	DiscountValue float64    `db:"discount_value"`
	MinOrderValue float64    `db:"min_order_value"`
	UsageLimit    *int       `db:"usage_limit"`
	UsedCount     int        `db:"used_count"`
	ValidFrom     *time.Time `db:"valid_from"`
	ValidTo       *time.Time `db:"valid_to"`
	Status        string     `db:"status"`
	CreatedBy     *string    `db:"created_by"`
	UpdatedBy     *string    `db:"updated_by"`
	DeletedBy     *string    `db:"deleted_by"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
	Version       int        `db:"version"`
}

const (
	DiscountTypePercentage = "percentage"
	DiscountTypeFixed      = "fixed"
)
