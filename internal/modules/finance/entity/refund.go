package entity

import "time"

// Refund mirrors the refunds table.
type Refund struct {
	ID           string     `db:"id"`
	TenantID     string     `db:"tenant_id"`
	PaymentID    string     `db:"payment_id"`
	Amount       float64    `db:"amount"`
	Reason       *string    `db:"reason"`
	RefundStatus string     `db:"refund_status"`
	ProcessedAt  *time.Time `db:"processed_at"`
	Status       string     `db:"status"`
	CreatedBy    *string    `db:"created_by"`
	UpdatedBy    *string    `db:"updated_by"`
	DeletedBy    *string    `db:"deleted_by"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	Version      int        `db:"version"`
}

const (
	RefundStatusPending   = "pending"
	RefundStatusProcessed = "processed"
	RefundStatusFailed    = "failed"
)
