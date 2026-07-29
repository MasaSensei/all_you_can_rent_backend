package service

import (
	"context"

	"rentos-backend/internal/modules/finance/dto/request"
	"rentos-backend/internal/modules/finance/dto/response"
)

// InvoiceService manages invoice lifecycle.
type InvoiceService interface {
	// CreateFromBooking generates an invoice from a confirmed booking,
	// applying the tenant's default tax rate to each line item.
	CreateFromBooking(ctx context.Context, tenantID, actorID string, req request.CreateInvoiceFromBooking) (*response.Invoice, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Invoice, error)
	List(ctx context.Context, tenantID string, filter request.ListInvoicesFilter) ([]response.Invoice, error)
}

// PaymentService manages payment recording and reconciliation.
type PaymentService interface {
	Record(ctx context.Context, tenantID, actorID string, req request.RecordPayment) (*response.Payment, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Payment, error)
	ListByInvoice(ctx context.Context, invoiceID, tenantID string) ([]response.Payment, error)
}

// RefundService manages refund issuance.
type RefundService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateRefund) (*response.Refund, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Refund, error)
}

// TaxService manages tax definitions.
type TaxService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateTax) (*response.Tax, error)
	List(ctx context.Context, tenantID string) ([]response.Tax, error)
}

// ExpenseService manages operational expenses.
type ExpenseService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateExpense) (*response.Expense, error)
	List(ctx context.Context, tenantID string, filter request.ListExpensesFilter) ([]response.Expense, error)
}

// BookingReader is the cross-module interface the finance module uses to
// read booking data without importing the booking package directly.
type BookingReader interface {
	GetItems(ctx context.Context, bookingID, tenantID string) ([]BookingItemSnapshot, error)
	GetTotals(ctx context.Context, bookingID, tenantID string) (BookingTotalsSnapshot, error)
}

// BookingItemSnapshot is the minimal booking item data finance needs.
type BookingItemSnapshot struct {
	ID        string
	AssetID   string
	Quantity  int
	UnitPrice float64
	LineTotal float64
}

// BookingTotalsSnapshot carries the booking-level aggregates finance uses
// to compute invoice totals.
type BookingTotalsSnapshot struct {
	Subtotal      float64
	DiscountTotal float64
}
