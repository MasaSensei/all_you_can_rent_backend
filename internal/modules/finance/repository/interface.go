package repository

import (
	"context"
	"errors"

	"rentos-backend/internal/modules/finance/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

// InvoiceRepository manages the invoices table.
type InvoiceRepository interface {
	Create(ctx context.Context, q database.Querier, inv *entity.Invoice) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Invoice, error)
	FindByBookingID(ctx context.Context, q database.Querier, bookingID, tenantID string) (*entity.Invoice, error)
	List(ctx context.Context, q database.Querier, tenantID string, customerID, status *string, limit, offset int) ([]entity.Invoice, error)
	// UpdatePaymentTotals recalculates amount_paid, amount_due, and invoice_status
	// based on the supplied amountPaid. Called after every payment or refund.
	UpdatePaymentTotals(ctx context.Context, q database.Querier, id, tenantID string, amountPaid float64) error
}

// InvoiceItemRepository manages the invoice_items table.
type InvoiceItemRepository interface {
	Create(ctx context.Context, q database.Querier, item *entity.InvoiceItem) error
	ListByInvoice(ctx context.Context, q database.Querier, invoiceID, tenantID string) ([]entity.InvoiceItem, error)
}

// PaymentRepository manages the payments table.
type PaymentRepository interface {
	Create(ctx context.Context, q database.Querier, p *entity.Payment) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Payment, error)
	ListByInvoice(ctx context.Context, q database.Querier, invoiceID, tenantID string) ([]entity.Payment, error)
	// SumByInvoice returns the total succeeded payment amount for an invoice.
	SumByInvoice(ctx context.Context, q database.Querier, invoiceID string) (float64, error)
}

// RefundRepository manages the refunds table.
type RefundRepository interface {
	Create(ctx context.Context, q database.Querier, r *entity.Refund) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Refund, error)
	// SumByPayment returns the total processed refund amount for a payment.
	SumByPayment(ctx context.Context, q database.Querier, paymentID string) (float64, error)
}

// TaxRepository manages the taxes table.
type TaxRepository interface {
	Create(ctx context.Context, q database.Querier, t *entity.Tax) error
	FindDefault(ctx context.Context, q database.Querier, tenantID string) (*entity.Tax, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Tax, error)
}

// ExpenseRepository manages the expenses table.
type ExpenseRepository interface {
	Create(ctx context.Context, q database.Querier, e *entity.Expense) error
	List(ctx context.Context, q database.Querier, tenantID string, assetID, category *string, limit, offset int) ([]entity.Expense, error)
}
