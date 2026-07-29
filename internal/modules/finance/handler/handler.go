package handler

import (
	"github.com/gofiber/fiber/v2"

	financereq "rentos-backend/internal/modules/finance/dto/request"
	"rentos-backend/internal/modules/finance/service"
	apiresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

const (
	ctxKeyTenantID = "tenant_id"
	ctxKeyUserID   = "user_id"
)

// Handler groups the finance module's HTTP handlers.
type Handler struct {
	invoices service.InvoiceService
	payments service.PaymentService
	refunds  service.RefundService
	taxes    service.TaxService
	expenses service.ExpenseService
	validate *validator.Validate
}

func New(
	invoices service.InvoiceService,
	payments service.PaymentService,
	refunds service.RefundService,
	taxes service.TaxService,
	expenses service.ExpenseService,
	v *validator.Validate,
) *Handler {
	return &Handler{
		invoices: invoices, payments: payments, refunds: refunds,
		taxes: taxes, expenses: expenses, validate: v,
	}
}

func tenantID(c *fiber.Ctx) string {
	if id, ok := c.Locals(ctxKeyTenantID).(string); ok {
		return id
	}
	return c.Get("X-Tenant-ID")
}

func userID(c *fiber.Ctx) string {
	id, _ := c.Locals(ctxKeyUserID).(string)
	return id
}

// ---- Invoices ----

func (h *Handler) CreateInvoiceFromBooking(c *fiber.Ctx) error {
	var req financereq.CreateInvoiceFromBooking
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	inv, err := h.invoices.CreateFromBooking(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, inv)
}

func (h *Handler) GetInvoice(c *fiber.Ctx) error {
	inv, err := h.invoices.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, inv)
}

func (h *Handler) ListInvoices(c *fiber.Ctx) error {
	filter := financereq.ListInvoicesFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("customer_id"); v != "" {
		filter.CustomerID = &v
	}
	if v := c.Query("invoice_status"); v != "" {
		filter.InvoiceStatus = &v
	}
	invoices, err := h.invoices.List(c.Context(), tenantID(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, invoices, fiber.Map{
		"page": filter.Page, "per_page": filter.PerPage,
	})
}

// ---- Payments ----

func (h *Handler) RecordPayment(c *fiber.Ctx) error {
	var req financereq.RecordPayment
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	p, err := h.payments.Record(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, p)
}

func (h *Handler) GetPayment(c *fiber.Ctx) error {
	p, err := h.payments.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, p)
}

func (h *Handler) ListPaymentsByInvoice(c *fiber.Ctx) error {
	payments, err := h.payments.ListByInvoice(c.Context(), c.Params("invoice_id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, payments)
}

// ---- Refunds ----

func (h *Handler) CreateRefund(c *fiber.Ctx) error {
	var req financereq.CreateRefund
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	ref, err := h.refunds.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, ref)
}

func (h *Handler) GetRefund(c *fiber.Ctx) error {
	ref, err := h.refunds.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, ref)
}

// ---- Taxes ----

func (h *Handler) CreateTax(c *fiber.Ctx) error {
	var req financereq.CreateTax
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	t, err := h.taxes.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, t)
}

func (h *Handler) ListTaxes(c *fiber.Ctx) error {
	taxes, err := h.taxes.List(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, taxes)
}

// ---- Expenses ----

func (h *Handler) CreateExpense(c *fiber.Ctx) error {
	var req financereq.CreateExpense
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	e, err := h.expenses.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, e)
}

func (h *Handler) ListExpenses(c *fiber.Ctx) error {
	filter := financereq.ListExpensesFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("asset_id"); v != "" {
		filter.AssetID = &v
	}
	if v := c.Query("category"); v != "" {
		filter.Category = &v
	}
	expenses, err := h.expenses.List(c.Context(), tenantID(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, expenses, fiber.Map{
		"page": filter.Page, "per_page": filter.PerPage,
	})
}
