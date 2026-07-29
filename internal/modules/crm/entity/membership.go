package entity

import "time"

// Membership mirrors the memberships table.
type Membership struct {
	ID               string     `db:"id"`
	TenantID         string     `db:"tenant_id"`
	CustomerID       string     `db:"customer_id"`
	PlanName         string     `db:"plan_name"`
	Tier             *string    `db:"tier"`
	StartDate        time.Time  `db:"start_date"`
	EndDate          *time.Time `db:"end_date"`
	Fee              float64    `db:"fee"`
	MembershipStatus string     `db:"membership_status"` // active, expired, cancelled
	Status           string     `db:"status"`
	CreatedBy        *string    `db:"created_by"`
	UpdatedBy        *string    `db:"updated_by"`
	DeletedBy        *string    `db:"deleted_by"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
	Version          int        `db:"version"`
}

const (
	MembershipStatusActive    = "active"
	MembershipStatusExpired   = "expired"
	MembershipStatusCancelled = "cancelled"
)
