package response

import "time"

// InvoiceItem is the API-facing shape of entity.InvoiceItem.
type InvoiceItem struct {
	ID            string  `json:"id"`
	BookingItemID *string `json:"booking_item_id,omitempty"`
	Description   string  `json:"description"`
	Quantity      int     `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	TaxAmount     float64 `json:"tax_amount"`
	LineTotal     float64 `json:"line_total"`
}

// Invoice is the full API-facing shape of an invoice.
type Invoice struct {
	ID            string        `json:"id"`
	TenantID      string        `json:"tenant_id"`
	CustomerID    string        `json:"customer_id"`
	BookingID     *string       `json:"booking_id,omitempty"`
	InvoiceNumber string        `json:"invoice_number"`
	IssueDate     time.Time     `json:"issue_date"`
	DueDate       time.Time     `json:"due_date"`
	Subtotal      float64       `json:"subtotal"`
	TaxTotal      float64       `json:"tax_total"`
	DiscountTotal float64       `json:"discount_total"`
	TotalAmount   float64       `json:"total_amount"`
	AmountPaid    float64       `json:"amount_paid"`
	AmountDue     float64       `json:"amount_due"`
	InvoiceStatus string        `json:"invoice_status"`
	Items         []InvoiceItem `json:"items,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// Payment is the API-facing shape of entity.Payment.
type Payment struct {
	ID                   string     `json:"id"`
	TenantID             string     `json:"tenant_id"`
	InvoiceID            string     `json:"invoice_id"`
	CustomerID           string     `json:"customer_id"`
	PaymentMethod        string     `json:"payment_method"`
	TransactionReference *string    `json:"transaction_reference,omitempty"`
	Amount               float64    `json:"amount"`
	Currency             string     `json:"currency"`
	PaidAt               *time.Time `json:"paid_at,omitempty"`
	PaymentStatus        string     `json:"payment_status"`
	CreatedAt            time.Time  `json:"created_at"`
}

// Refund is the API-facing shape of entity.Refund.
type Refund struct {
	ID           string     `json:"id"`
	PaymentID    string     `json:"payment_id"`
	Amount       float64    `json:"amount"`
	Reason       *string    `json:"reason,omitempty"`
	RefundStatus string     `json:"refund_status"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// Tax is the API-facing shape of entity.Tax.
type Tax struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Rate      float64   `json:"rate"`
	TaxType   string    `json:"tax_type"`
	IsDefault bool      `json:"is_default"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Expense is the API-facing shape of entity.Expense.
type Expense struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	AssetID     *string   `json:"asset_id,omitempty"`
	Category    string    `json:"category"`
	Amount      float64   `json:"amount"`
	ExpenseDate time.Time `json:"expense_date"`
	Description *string   `json:"description,omitempty"`
	Vendor      *string   `json:"vendor,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
