package entity

import "time"

// Payment mirrors the payments table.
type Payment struct {
	ID                   string     `db:"id"`
	TenantID             string     `db:"tenant_id"`
	InvoiceID            string     `db:"invoice_id"`
	CustomerID           string     `db:"customer_id"`
	PaymentMethod        string     `db:"payment_method"`
	TransactionReference *string    `db:"transaction_reference"`
	Amount               float64    `db:"amount"`
	Currency             string     `db:"currency"`
	PaidAt               *time.Time `db:"paid_at"`
	PaymentStatus        string     `db:"payment_status"`
	Status               string     `db:"status"`
	CreatedBy            *string    `db:"created_by"`
	UpdatedBy            *string    `db:"updated_by"`
	DeletedBy            *string    `db:"deleted_by"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
	DeletedAt            *time.Time `db:"deleted_at"`
	Version              int        `db:"version"`
}

const (
	PaymentStatusPending   = "pending"
	PaymentStatusSucceeded = "succeeded"
	PaymentStatusFailed    = "failed"
	PaymentStatusRefunded  = "refunded"
)
