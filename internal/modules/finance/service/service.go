package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/finance/dto/request"
	"rentos-backend/internal/modules/finance/dto/response"
	"rentos-backend/internal/modules/finance/entity"
	"rentos-backend/internal/modules/finance/repository"
	pkgresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/transaction"
)

// ============================================================
// invoiceService
// ============================================================

type invoiceService struct {
	db      *sqlx.DB
	invoices repository.InvoiceRepository
	items    repository.InvoiceItemRepository
	taxes    repository.TaxRepository
	booking  BookingReader
}

func NewInvoiceService(
	db *sqlx.DB,
	invoices repository.InvoiceRepository,
	items repository.InvoiceItemRepository,
	taxes repository.TaxRepository,
	booking BookingReader,
) InvoiceService {
	return &invoiceService{db: db, invoices: invoices, items: items, taxes: taxes, booking: booking}
}

func (s *invoiceService) CreateFromBooking(ctx context.Context, tenantID, actorID string, req request.CreateInvoiceFromBooking) (*response.Invoice, error) {
	// Guard: invoice must not already exist for this booking.
	if existing, err := s.invoices.FindByBookingID(ctx, s.db, req.BookingID, tenantID); err == nil && existing != nil {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "an invoice already exists for this booking")
	}

	// Read booking snapshot via cross-module interface.
	bookingItems, err := s.booking.GetItems(ctx, req.BookingID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("finance: read booking items: %w", err)
	}
	totals, err := s.booking.GetTotals(ctx, req.BookingID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("finance: read booking totals: %w", err)
	}

	// Resolve default tax rate (0 if none configured).
	var taxRate float64
	if tax, err := s.taxes.FindDefault(ctx, s.db, tenantID); err == nil {
		taxRate = tax.Rate / 100
	}

	// Build invoice items with per-line tax.
	type lineItem struct {
		entity    entity.InvoiceItem
		lineTotal float64
	}
	lines := make([]lineItem, 0, len(bookingItems))
	var subtotal float64

	for _, bi := range bookingItems {
		taxAmount := round2(bi.LineTotal * taxRate)
		lines = append(lines, lineItem{
			entity: entity.InvoiceItem{
				ID:            uuid.NewString(),
				TenantID:      tenantID,
				BookingItemID: &bi.ID,
				Description:   fmt.Sprintf("Asset %s × %d", bi.AssetID, bi.Quantity),
				Quantity:      bi.Quantity,
				UnitPrice:     bi.UnitPrice,
				TaxAmount:     taxAmount,
				LineTotal:     round2(bi.LineTotal + taxAmount),
			},
			lineTotal: round2(bi.LineTotal + taxAmount),
		})
		subtotal += bi.LineTotal
	}

	taxTotal := round2(subtotal * taxRate)
	totalAmount := round2(subtotal + taxTotal - totals.DiscountTotal)

	inv := &entity.Invoice{
		ID:            uuid.NewString(),
		TenantID:      tenantID,
		CustomerID:    req.CustomerID,
		BookingID:     &req.BookingID,
		InvoiceNumber: generateInvoiceNumber(),
		DueDate:       req.DueDate,
		Subtotal:      round2(subtotal),
		TaxTotal:      taxTotal,
		DiscountTotal: round2(totals.DiscountTotal),
		TotalAmount:   totalAmount,
		CreatedBy:     &actorID,
	}

	if err := transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if err := s.invoices.Create(ctx, tx, inv); err != nil {
			return err
		}
		for i := range lines {
			lines[i].entity.InvoiceID = inv.ID
			if err := s.items.Create(ctx, tx, &lines[i].entity); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, inv.ID, tenantID)
}

func (s *invoiceService) GetByID(ctx context.Context, id, tenantID string) (*response.Invoice, error) {
	inv, err := s.invoices.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "invoice not found")
		}
		return nil, err
	}
	return s.hydrate(ctx, inv)
}

func (s *invoiceService) List(ctx context.Context, tenantID string, filter request.ListInvoicesFilter) ([]response.Invoice, error) {
	perPage, page := normPage(filter.PerPage, filter.Page)
	invoices, err := s.invoices.List(ctx, s.db, tenantID, filter.CustomerID, filter.InvoiceStatus, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.Invoice, 0, len(invoices))
	for _, inv := range invoices {
		r, err := s.hydrate(ctx, &inv)
		if err != nil {
			return nil, err
		}
		out = append(out, *r)
	}
	return out, nil
}

func (s *invoiceService) hydrate(ctx context.Context, inv *entity.Invoice) (*response.Invoice, error) {
	items, err := s.items.ListByInvoice(ctx, s.db, inv.ID, inv.TenantID)
	if err != nil {
		return nil, err
	}
	itemResp := make([]response.InvoiceItem, 0, len(items))
	for _, it := range items {
		itemResp = append(itemResp, response.InvoiceItem{
			ID: it.ID, BookingItemID: it.BookingItemID,
			Description: it.Description, Quantity: it.Quantity,
			UnitPrice: it.UnitPrice, TaxAmount: it.TaxAmount, LineTotal: it.LineTotal,
		})
	}
	return &response.Invoice{
		ID: inv.ID, TenantID: inv.TenantID, CustomerID: inv.CustomerID,
		BookingID: inv.BookingID, InvoiceNumber: inv.InvoiceNumber,
		IssueDate: inv.IssueDate, DueDate: inv.DueDate,
		Subtotal: inv.Subtotal, TaxTotal: inv.TaxTotal, DiscountTotal: inv.DiscountTotal,
		TotalAmount: inv.TotalAmount, AmountPaid: inv.AmountPaid, AmountDue: inv.AmountDue,
		InvoiceStatus: inv.InvoiceStatus, Items: itemResp,
		CreatedAt: inv.CreatedAt, UpdatedAt: inv.UpdatedAt,
	}, nil
}

// ============================================================
// paymentService
// ============================================================

type paymentService struct {
	db       *sqlx.DB
	payments repository.PaymentRepository
	invoices repository.InvoiceRepository
}

func NewPaymentService(
	db *sqlx.DB,
	payments repository.PaymentRepository,
	invoices repository.InvoiceRepository,
) PaymentService {
	return &paymentService{db: db, payments: payments, invoices: invoices}
}

func (s *paymentService) Record(ctx context.Context, tenantID, actorID string, req request.RecordPayment) (*response.Payment, error) {
	inv, err := s.invoices.FindByID(ctx, s.db, req.InvoiceID, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "invoice not found")
		}
		return nil, err
	}
	if inv.InvoiceStatus == entity.InvoiceStatusPaid {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "invoice is already fully paid")
	}

	p := &entity.Payment{
		ID:                   uuid.NewString(),
		TenantID:             tenantID,
		InvoiceID:            req.InvoiceID,
		CustomerID:           req.CustomerID,
		PaymentMethod:        req.PaymentMethod,
		TransactionReference: req.TransactionReference,
		Amount:               req.Amount,
		Currency:             req.Currency,
		CreatedBy:            &actorID,
	}

	if err := transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if err := s.payments.Create(ctx, tx, p); err != nil {
			return err
		}
		// Recalculate total paid and update invoice status atomically.
		totalPaid, err := s.payments.SumByInvoice(ctx, tx, req.InvoiceID)
		if err != nil {
			return err
		}
		return s.invoices.UpdatePaymentTotals(ctx, tx, req.InvoiceID, tenantID, totalPaid)
	}); err != nil {
		return nil, err
	}

	return &response.Payment{
		ID: p.ID, TenantID: p.TenantID, InvoiceID: p.InvoiceID,
		CustomerID: p.CustomerID, PaymentMethod: p.PaymentMethod,
		TransactionReference: p.TransactionReference,
		Amount: p.Amount, Currency: p.Currency,
		PaymentStatus: entity.PaymentStatusSucceeded,
		CreatedAt: p.CreatedAt,
	}, nil
}

func (s *paymentService) GetByID(ctx context.Context, id, tenantID string) (*response.Payment, error) {
	p, err := s.payments.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "payment not found")
		}
		return nil, err
	}
	return &response.Payment{
		ID: p.ID, TenantID: p.TenantID, InvoiceID: p.InvoiceID,
		CustomerID: p.CustomerID, PaymentMethod: p.PaymentMethod,
		TransactionReference: p.TransactionReference,
		Amount: p.Amount, Currency: p.Currency,
		PaidAt: p.PaidAt, PaymentStatus: p.PaymentStatus,
		CreatedAt: p.CreatedAt,
	}, nil
}

func (s *paymentService) ListByInvoice(ctx context.Context, invoiceID, tenantID string) ([]response.Payment, error) {
	payments, err := s.payments.ListByInvoice(ctx, s.db, invoiceID, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.Payment, 0, len(payments))
	for _, p := range payments {
		out = append(out, response.Payment{
			ID: p.ID, TenantID: p.TenantID, InvoiceID: p.InvoiceID,
			CustomerID: p.CustomerID, PaymentMethod: p.PaymentMethod,
			TransactionReference: p.TransactionReference,
			Amount: p.Amount, Currency: p.Currency,
			PaidAt: p.PaidAt, PaymentStatus: p.PaymentStatus,
			CreatedAt: p.CreatedAt,
		})
	}
	return out, nil
}

// ============================================================
// refundService
// ============================================================

type refundService struct {
	db       *sqlx.DB
	refunds  repository.RefundRepository
	payments repository.PaymentRepository
	invoices repository.InvoiceRepository
}

func NewRefundService(
	db *sqlx.DB,
	refunds repository.RefundRepository,
	payments repository.PaymentRepository,
	invoices repository.InvoiceRepository,
) RefundService {
	return &refundService{db: db, refunds: refunds, payments: payments, invoices: invoices}
}

func (s *refundService) Create(ctx context.Context, tenantID, actorID string, req request.CreateRefund) (*response.Refund, error) {
	p, err := s.payments.FindByID(ctx, s.db, req.PaymentID, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "payment not found")
		}
		return nil, err
	}

	// Cannot refund more than the net paid amount (payment − existing refunds).
	alreadyRefunded, err := s.refunds.SumByPayment(ctx, s.db, req.PaymentID)
	if err != nil {
		return nil, err
	}
	maxRefundable := p.Amount - alreadyRefunded
	if req.Amount > maxRefundable {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict,
			fmt.Sprintf("refund amount %.2f exceeds maximum refundable amount %.2f", req.Amount, maxRefundable))
	}

	ref := &entity.Refund{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		PaymentID: req.PaymentID,
		Amount:    req.Amount,
		Reason:    req.Reason,
		CreatedBy: &actorID,
	}

	if err := transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if err := s.refunds.Create(ctx, tx, ref); err != nil {
			return err
		}
		// Recompute invoice amount_paid after the refund.
		totalPaid, err := s.payments.SumByInvoice(ctx, tx, p.InvoiceID)
		if err != nil {
			return err
		}
		totalRefunded, err := s.refunds.SumByPayment(ctx, tx, req.PaymentID)
		if err != nil {
			return err
		}
		return s.invoices.UpdatePaymentTotals(ctx, tx, p.InvoiceID, tenantID, totalPaid-totalRefunded)
	}); err != nil {
		return nil, err
	}

	return &response.Refund{
		ID: ref.ID, PaymentID: ref.PaymentID, Amount: ref.Amount,
		Reason: ref.Reason, RefundStatus: entity.RefundStatusProcessed,
		ProcessedAt: ref.ProcessedAt, CreatedAt: ref.CreatedAt,
	}, nil
}

func (s *refundService) GetByID(ctx context.Context, id, tenantID string) (*response.Refund, error) {
	ref, err := s.refunds.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "refund not found")
		}
		return nil, err
	}
	return &response.Refund{
		ID: ref.ID, PaymentID: ref.PaymentID, Amount: ref.Amount,
		Reason: ref.Reason, RefundStatus: ref.RefundStatus,
		ProcessedAt: ref.ProcessedAt, CreatedAt: ref.CreatedAt,
	}, nil
}

// ============================================================
// taxService
// ============================================================

type taxService struct {
	db   *sqlx.DB
	repo repository.TaxRepository
}

func NewTaxService(db *sqlx.DB, repo repository.TaxRepository) TaxService {
	return &taxService{db: db, repo: repo}
}

func (s *taxService) Create(ctx context.Context, tenantID, actorID string, req request.CreateTax) (*response.Tax, error) {
	t := &entity.Tax{
		ID:        uuid.NewString(),
		TenantID:  tenantID,
		Name:      req.Name,
		Rate:      req.Rate,
		TaxType:   req.TaxType,
		IsDefault: req.IsDefault,
		CreatedBy: &actorID,
	}
	if err := s.repo.Create(ctx, s.db, t); err != nil {
		return nil, err
	}
	return &response.Tax{
		ID: t.ID, TenantID: t.TenantID, Name: t.Name,
		Rate: t.Rate, TaxType: t.TaxType, IsDefault: t.IsDefault,
		Status: "active", CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}, nil
}

func (s *taxService) List(ctx context.Context, tenantID string) ([]response.Tax, error) {
	taxes, err := s.repo.List(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.Tax, 0, len(taxes))
	for _, t := range taxes {
		out = append(out, response.Tax{
			ID: t.ID, TenantID: t.TenantID, Name: t.Name,
			Rate: t.Rate, TaxType: t.TaxType, IsDefault: t.IsDefault,
			Status: t.Status, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
		})
	}
	return out, nil
}

// ============================================================
// expenseService
// ============================================================

type expenseService struct {
	db   *sqlx.DB
	repo repository.ExpenseRepository
}

func NewExpenseService(db *sqlx.DB, repo repository.ExpenseRepository) ExpenseService {
	return &expenseService{db: db, repo: repo}
}

func (s *expenseService) Create(ctx context.Context, tenantID, actorID string, req request.CreateExpense) (*response.Expense, error) {
	e := &entity.Expense{
		ID:          uuid.NewString(),
		TenantID:    tenantID,
		AssetID:     req.AssetID,
		Category:    req.Category,
		Amount:      req.Amount,
		ExpenseDate: req.ExpenseDate,
		Description: req.Description,
		Vendor:      req.Vendor,
		CreatedBy:   &actorID,
	}
	if err := s.repo.Create(ctx, s.db, e); err != nil {
		return nil, err
	}
	return &response.Expense{
		ID: e.ID, TenantID: e.TenantID, AssetID: e.AssetID,
		Category: e.Category, Amount: e.Amount, ExpenseDate: e.ExpenseDate,
		Description: e.Description, Vendor: e.Vendor,
		Status: "active", CreatedAt: e.CreatedAt,
	}, nil
}

func (s *expenseService) List(ctx context.Context, tenantID string, filter request.ListExpensesFilter) ([]response.Expense, error) {
	perPage, page := normPage(filter.PerPage, filter.Page)
	expenses, err := s.repo.List(ctx, s.db, tenantID, filter.AssetID, filter.Category, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.Expense, 0, len(expenses))
	for _, e := range expenses {
		out = append(out, response.Expense{
			ID: e.ID, TenantID: e.TenantID, AssetID: e.AssetID,
			Category: e.Category, Amount: e.Amount, ExpenseDate: e.ExpenseDate,
			Description: e.Description, Vendor: e.Vendor,
			Status: e.Status, CreatedAt: e.CreatedAt,
		})
	}
	return out, nil
}

// ============================================================
// helpers
// ============================================================

func generateInvoiceNumber() string {
	return fmt.Sprintf("INV-%s", uuid.NewString()[:8])
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func normPage(perPage, page int) (int, int) {
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	if page <= 0 {
		page = 1
	}
	return perPage, page
}
