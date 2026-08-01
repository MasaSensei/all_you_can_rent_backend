package entity

import "time"

// Invoice mirrors the invoices table.
type Invoice struct {
	ID            string     `db:"id"`
	TenantID      string     `db:"tenant_id"`
	CustomerID    string     `db:"customer_id"`
	BookingID     *string    `db:"booking_id"`
	InvoiceNumber string     `db:"invoice_number"`
	IssueDate     time.Time  `db:"issue_date"`
	DueDate       time.Time  `db:"due_date"`
	Subtotal      float64    `db:"subtotal"`
	TaxTotal      float64    `db:"tax_total"`
	DiscountTotal float64    `db:"discount_total"`
	TotalAmount   float64    `db:"total_amount"`
	AmountPaid    float64    `db:"amount_paid"`
	AmountDue     float64    `db:"amount_due"`
	InvoiceStatus string     `db:"invoice_status"`
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
	InvoiceStatusUnpaid  = "unpaid"
	InvoiceStatusPartial = "partial"
	InvoiceStatusPaid    = "paid"
	InvoiceStatusVoided  = "voided"
)
