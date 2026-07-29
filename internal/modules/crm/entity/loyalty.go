package entity

import "time"

// LoyaltyProgram mirrors the loyalty_programs table.
type LoyaltyProgram struct {
	ID                 string     `db:"id"`
	TenantID           string     `db:"tenant_id"`
	Name               string     `db:"name"`
	Description        *string    `db:"description"`
	PointsPerCurrency  float64    `db:"points_per_currency"`
	RedemptionRate     float64    `db:"redemption_rate"`
	Status             string     `db:"status"`
	CreatedBy          *string    `db:"created_by"`
	UpdatedBy          *string    `db:"updated_by"`
	DeletedBy          *string    `db:"deleted_by"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
	DeletedAt          *time.Time `db:"deleted_at"`
	Version            int        `db:"version"`
}

// LoyaltyTransaction mirrors the loyalty_transactions table.
type LoyaltyTransaction struct {
	ID               string     `db:"id"`
	TenantID         string     `db:"tenant_id"`
	LoyaltyProgramID string     `db:"loyalty_program_id"`
	CustomerID       string     `db:"customer_id"`
	BookingID        *string    `db:"booking_id"`
	Points           int        `db:"points"`
	TransactionType  string     `db:"transaction_type"` // earn, redeem, expire, adjust
	Description      *string    `db:"description"`
	Status           string     `db:"status"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
	Version          int        `db:"version"`
}

const (
	LoyaltyTxTypeEarn   = "earn"
	LoyaltyTxTypeRedeem = "redeem"
	LoyaltyTxTypeExpire = "expire"
	LoyaltyTxTypeAdjust = "adjust"
)
