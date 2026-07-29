package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos-backend/internal/modules/finance/entity"
	"rentos-backend/internal/modules/finance/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// invoiceRepository
// ============================================================

type invoiceRepository struct {
	qCreate              string
	qFindByID            string
	qFindByBookingID     string
	qList                string
	qUpdatePaymentTotals string
}

func NewInvoiceRepository(qCreate, qFindByID, qFindByBookingID, qList, qUpdatePaymentTotals string) repository.InvoiceRepository {
	return &invoiceRepository{
		qCreate: qCreate, qFindByID: qFindByID,
		qFindByBookingID: qFindByBookingID,
		qList: qList, qUpdatePaymentTotals: qUpdatePaymentTotals,
	}
}

func (r *invoiceRepository) Create(ctx context.Context, q database.Querier, inv *entity.Invoice) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		inv.ID, inv.TenantID, inv.CustomerID, inv.BookingID, inv.InvoiceNumber,
		inv.DueDate, inv.Subtotal, inv.TaxTotal, inv.DiscountTotal, inv.TotalAmount,
		inv.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("invoiceRepository.Create: %w", err)
	}
	return nil
}

func (r *invoiceRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Invoice, error) {
	var inv entity.Invoice
	if err := q.GetContext(ctx, &inv, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("invoiceRepository.FindByID: %w", err)
	}
	return &inv, nil
}

func (r *invoiceRepository) FindByBookingID(ctx context.Context, q database.Querier, bookingID, tenantID string) (*entity.Invoice, error) {
	var inv entity.Invoice
	if err := q.GetContext(ctx, &inv, r.qFindByBookingID, bookingID, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("invoiceRepository.FindByBookingID: %w", err)
	}
	return &inv, nil
}

func (r *invoiceRepository) List(ctx context.Context, q database.Querier, tenantID string, customerID, status *string, limit, offset int) ([]entity.Invoice, error) {
	var out []entity.Invoice
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, customerID, status, limit, offset); err != nil {
		return nil, fmt.Errorf("invoiceRepository.List: %w", err)
	}
	return out, nil
}

func (r *invoiceRepository) UpdatePaymentTotals(ctx context.Context, q database.Querier, id, tenantID string, amountPaid float64) error {
	_, err := q.ExecContext(ctx, r.qUpdatePaymentTotals, id, tenantID, amountPaid)
	if err != nil {
		return fmt.Errorf("invoiceRepository.UpdatePaymentTotals: %w", err)
	}
	return nil
}

// ============================================================
// invoiceItemRepository
// ============================================================

type invoiceItemRepository struct {
	qCreate string
	qList   string
}

func NewInvoiceItemRepository(qCreate, qList string) repository.InvoiceItemRepository {
	return &invoiceItemRepository{qCreate: qCreate, qList: qList}
}

func (r *invoiceItemRepository) Create(ctx context.Context, q database.Querier, item *entity.InvoiceItem) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		item.ID, item.TenantID, item.InvoiceID, item.BookingItemID,
		item.Description, item.Quantity, item.UnitPrice, item.TaxAmount, item.LineTotal,
	)
	if err != nil {
		return fmt.Errorf("invoiceItemRepository.Create: %w", err)
	}
	return nil
}

func (r *invoiceItemRepository) ListByInvoice(ctx context.Context, q database.Querier, invoiceID, tenantID string) ([]entity.InvoiceItem, error) {
	var out []entity.InvoiceItem
	if err := q.SelectContext(ctx, &out, r.qList, invoiceID, tenantID); err != nil {
		return nil, fmt.Errorf("invoiceItemRepository.ListByInvoice: %w", err)
	}
	return out, nil
}

// ============================================================
// paymentRepository
// ============================================================

type paymentRepository struct {
	qCreate        string
	qFindByID      string
	qListByInvoice string
	qSumByInvoice  string
}

func NewPaymentRepository(qCreate, qFindByID, qListByInvoice, qSumByInvoice string) repository.PaymentRepository {
	return &paymentRepository{
		qCreate: qCreate, qFindByID: qFindByID,
		qListByInvoice: qListByInvoice, qSumByInvoice: qSumByInvoice,
	}
}

func (r *paymentRepository) Create(ctx context.Context, q database.Querier, p *entity.Payment) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		p.ID, p.TenantID, p.InvoiceID, p.CustomerID, p.PaymentMethod,
		p.TransactionReference, p.Amount, p.Currency, p.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("paymentRepository.Create: %w", err)
	}
	return nil
}

func (r *paymentRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Payment, error) {
	var p entity.Payment
	if err := q.GetContext(ctx, &p, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("paymentRepository.FindByID: %w", err)
	}
	return &p, nil
}

func (r *paymentRepository) ListByInvoice(ctx context.Context, q database.Querier, invoiceID, tenantID string) ([]entity.Payment, error) {
	var out []entity.Payment
	if err := q.SelectContext(ctx, &out, r.qListByInvoice, invoiceID, tenantID); err != nil {
		return nil, fmt.Errorf("paymentRepository.ListByInvoice: %w", err)
	}
	return out, nil
}

func (r *paymentRepository) SumByInvoice(ctx context.Context, q database.Querier, invoiceID string) (float64, error) {
	var total float64
	if err := q.GetContext(ctx, &total, r.qSumByInvoice, invoiceID); err != nil {
		return 0, fmt.Errorf("paymentRepository.SumByInvoice: %w", err)
	}
	return total, nil
}
