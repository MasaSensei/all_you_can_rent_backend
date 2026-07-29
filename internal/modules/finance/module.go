package finance

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/finance/handler"
	"rentos-backend/internal/modules/finance/repository/postgres"
	"rentos-backend/internal/modules/finance/routes"
	"rentos-backend/internal/modules/finance/service"
	"rentos-backend/pkg/database"
)

// Module holds the finance module's wired handler.
type Module struct {
	handler *handler.Handler
}

// BookingRepositoryAdapter allows the finance module to read booking data
// from the database directly without importing the booking service package.
// It queries booking_items and bookings using embedded SQL, satisfying
// service.BookingReader via the module-level query() embed loader.
type BookingRepositoryAdapter struct {
	db             *sqlx.DB
	qGetItems      string
	qGetTotals     string
}

func (a *BookingRepositoryAdapter) GetItems(ctx context.Context, bookingID, tenantID string) ([]service.BookingItemSnapshot, error) {
	type row struct {
		ID        string  `db:"id"`
		AssetID   string  `db:"asset_id"`
		Quantity  int     `db:"quantity"`
		UnitPrice float64 `db:"unit_price"`
		LineTotal float64 `db:"line_total"`
	}
	var rows []row
	if err := a.db.SelectContext(ctx, &rows, a.qGetItems, bookingID, tenantID); err != nil {
		return nil, fmt.Errorf("BookingRepositoryAdapter.GetItems: %w", err)
	}
	out := make([]service.BookingItemSnapshot, 0, len(rows))
	for _, r := range rows {
		out = append(out, service.BookingItemSnapshot{
			ID: r.ID, AssetID: r.AssetID,
			Quantity: r.Quantity, UnitPrice: r.UnitPrice, LineTotal: r.LineTotal,
		})
	}
	return out, nil
}

func (a *BookingRepositoryAdapter) GetTotals(ctx context.Context, bookingID, tenantID string) (service.BookingTotalsSnapshot, error) {
	var snap struct {
		Subtotal      float64 `db:"subtotal"`
		DiscountTotal float64 `db:"discount_total"`
	}
	if err := a.db.GetContext(ctx, &snap, a.qGetTotals, bookingID, tenantID); err != nil {
		return service.BookingTotalsSnapshot{}, fmt.Errorf("BookingRepositoryAdapter.GetTotals: %w", err)
	}
	return service.BookingTotalsSnapshot{
		Subtotal:      snap.Subtotal,
		DiscountTotal: snap.DiscountTotal,
	}, nil
}

// Ensure adapter satisfies the interface at compile time.
var _ service.BookingReader = (*BookingRepositoryAdapter)(nil)
var _ database.Querier = (*sqlx.DB)(nil)

// New builds the finance module.
func New(c *bootstrap.Container) *Module {
	invoiceRepo := postgres.NewInvoiceRepository(
		query("create_invoice.sql"),
		query("find_invoice_by_id.sql"),
		query("find_invoice_by_booking_id.sql"),
		query("list_invoices.sql"),
		query("update_invoice_payment_totals.sql"),
	)
	itemRepo := postgres.NewInvoiceItemRepository(
		query("create_invoice_item.sql"),
		query("list_invoice_items.sql"),
	)
	paymentRepo := postgres.NewPaymentRepository(
		query("create_payment.sql"),
		query("find_payment_by_id.sql"),
		query("list_payments_by_invoice.sql"),
		query("sum_payments_by_invoice.sql"),
	)
	refundRepo := postgres.NewRefundRepository(
		query("create_refund.sql"),
		query("find_refund_by_id.sql"),
		query("sum_refunds_by_payment.sql"),
	)
	taxRepo := postgres.NewTaxRepository(
		query("create_tax.sql"),
		query("find_default_tax.sql"),
		query("list_taxes.sql"),
	)
	expenseRepo := postgres.NewExpenseRepository(
		query("create_expense.sql"),
		query("list_expenses.sql"),
	)

	bookingReader := &BookingRepositoryAdapter{
		db:         c.DB,
		qGetItems:  query("get_booking_items_for_invoice.sql"),
		qGetTotals: query("get_booking_totals_for_invoice.sql"),
	}

	invoiceSvc := service.NewInvoiceService(c.DB, invoiceRepo, itemRepo, taxRepo, bookingReader)
	paymentSvc := service.NewPaymentService(c.DB, paymentRepo, invoiceRepo)
	refundSvc := service.NewRefundService(c.DB, refundRepo, paymentRepo, invoiceRepo)
	taxSvc := service.NewTaxService(c.DB, taxRepo)
	expenseSvc := service.NewExpenseService(c.DB, expenseRepo)

	h := handler.New(invoiceSvc, paymentSvc, refundSvc, taxSvc, expenseSvc, c.Validator)
	return &Module{handler: h}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}
