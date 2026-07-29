package entity

import "time"

// Promotion mirrors the promotions table.
type Promotion struct {
	ID            string     `db:"id"`
	TenantID      string     `db:"tenant_id"`
	Name          string     `db:"name"`
	Description   *string    `db:"description"`
	PromotionType string     `db:"promotion_type"` // percentage, fixed
	Value         float64    `db:"value"`
	StartDate     *time.Time `db:"start_date"`
	EndDate       *time.Time `db:"end_date"`
	Status        string     `db:"status"`
	CreatedBy     *string    `db:"created_by"`
	UpdatedBy     *string    `db:"updated_by"`
	DeletedBy     *string    `db:"deleted_by"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
	Version       int        `db:"version"`
}
