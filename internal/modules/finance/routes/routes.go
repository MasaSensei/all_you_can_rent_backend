package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/finance/handler"
)

// Register mounts all finance routes under /api/v1.
func Register(router fiber.Router, h *handler.Handler) {
	invoices := router.Group("/invoices")
	invoices.Post("/", h.CreateInvoiceFromBooking)
	invoices.Get("/", h.ListInvoices)
	invoices.Get("/:id", h.GetInvoice)
	invoices.Get("/:invoice_id/payments", h.ListPaymentsByInvoice)

	payments := router.Group("/payments")
	payments.Post("/", h.RecordPayment)
	payments.Get("/:id", h.GetPayment)

	refunds := router.Group("/refunds")
	refunds.Post("/", h.CreateRefund)
	refunds.Get("/:id", h.GetRefund)

	taxes := router.Group("/taxes")
	taxes.Post("/", h.CreateTax)
	taxes.Get("/", h.ListTaxes)

	expenses := router.Group("/expenses")
	expenses.Post("/", h.CreateExpense)
	expenses.Get("/", h.ListExpenses)
}
