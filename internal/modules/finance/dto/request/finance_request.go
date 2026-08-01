package request

import "time"

// CreateInvoiceFromBooking generates an invoice from an existing booking.
type CreateInvoiceFromBooking struct {
	BookingID  string    `json:"booking_id" validate:"required,uuid"`
	CustomerID string    `json:"customer_id" validate:"required,uuid"`
	DueDate    time.Time `json:"due_date" validate:"required"`
}

// RecordPayment records a payment against an invoice.
type RecordPayment struct {
	InvoiceID            string  `json:"invoice_id" validate:"required,uuid"`
	CustomerID           string  `json:"customer_id" validate:"required,uuid"`
	PaymentMethod        string  `json:"payment_method" validate:"required,oneof=cash card bank_transfer online"`
	TransactionReference *string `json:"transaction_reference" validate:"omitempty,max=150"`
	Amount               float64 `json:"amount" validate:"required,min=0.01"`
	Currency             string  `json:"currency" validate:"required,len=3"`
}

// CreateRefund issues a refund against a payment.
type CreateRefund struct {
	PaymentID string  `json:"payment_id" validate:"required,uuid"`
	Amount    float64 `json:"amount" validate:"required,min=0.01"`
	Reason    *string `json:"reason" validate:"omitempty,max=500"`
}

// CreateTax creates a new tax definition.
type CreateTax struct {
	Name      string  `json:"name" validate:"required,max=100"`
	Rate      float64 `json:"rate" validate:"required,min=0,max=100"`
	TaxType   string  `json:"tax_type" validate:"required,oneof=percentage fixed"`
	IsDefault bool    `json:"is_default"`
}

// CreateExpense records an operational expense.
type CreateExpense struct {
	AssetID     *string   `json:"asset_id" validate:"omitempty,uuid"`
	Category    string    `json:"category" validate:"required,max=100"`
	Amount      float64   `json:"amount" validate:"required,min=0.01"`
	ExpenseDate time.Time `json:"expense_date" validate:"required"`
	Description *string   `json:"description" validate:"omitempty,max=500"`
	Vendor      *string   `json:"vendor" validate:"omitempty,max=255"`
}

// ListInvoicesFilter holds whitelisted query-param filters.
type ListInvoicesFilter struct {
	CustomerID    *string
	InvoiceStatus *string
	Page          int
	PerPage       int
}

// ListExpensesFilter holds whitelisted query-param filters.
type ListExpensesFilter struct {
	AssetID  *string
	Category *string
	Page     int
	PerPage  int
}
