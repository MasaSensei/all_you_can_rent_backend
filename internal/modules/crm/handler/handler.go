package handler

import (
	"github.com/gofiber/fiber/v2"

	crmreq "rentos-backend/internal/modules/crm/dto/request"
	"rentos-backend/internal/modules/crm/service"
	apiresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

const (
	ctxKeyTenantID = "tenant_id"
	ctxKeyUserID   = "user_id"
)

// Handler groups all CRM HTTP handlers.
type Handler struct {
	customers   service.CustomerService
	memberships service.MembershipService
	loyalty     service.LoyaltyService
	validate    *validator.Validate
}

func New(
	customers service.CustomerService,
	memberships service.MembershipService,
	loyalty service.LoyaltyService,
	v *validator.Validate,
) *Handler {
	return &Handler{customers: customers, memberships: memberships, loyalty: loyalty, validate: v}
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

// ---- Customers ----

func (h *Handler) CreateCustomer(c *fiber.Ctx) error {
	var req crmreq.CreateCustomer
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	customer, err := h.customers.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, customer)
}

func (h *Handler) GetCustomer(c *fiber.Ctx) error {
	customer, err := h.customers.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, customer)
}

func (h *Handler) ListCustomers(c *fiber.Ctx) error {
	filter := crmreq.ListCustomersFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("customer_type"); v != "" {
		filter.CustomerType = &v
	}
	if v := c.Query("search"); v != "" {
		filter.Search = &v
	}
	customers, err := h.customers.List(c.Context(), tenantID(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, customers, fiber.Map{
		"page": filter.Page, "per_page": filter.PerPage,
	})
}

func (h *Handler) UpdateCustomer(c *fiber.Ctx) error {
	var req crmreq.UpdateCustomer
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	customer, err := h.customers.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, customer)
}

func (h *Handler) DeleteCustomer(c *fiber.Ctx) error {
	if err := h.customers.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Customer Addresses ----

func (h *Handler) AddAddress(c *fiber.Ctx) error {
	var req crmreq.AddAddress
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	addr, err := h.customers.AddAddress(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, addr)
}

func (h *Handler) DeleteAddress(c *fiber.Ctx) error {
	if err := h.customers.DeleteAddress(c.Context(), c.Params("addr_id"), c.Params("id"), tenantID(c), userID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Memberships ----

func (h *Handler) CreateMembership(c *fiber.Ctx) error {
	var req crmreq.CreateMembership
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	m, err := h.memberships.Create(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, m)
}

func (h *Handler) ListMemberships(c *fiber.Ctx) error {
	members, err := h.memberships.ListByCustomer(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, members)
}

// ---- Loyalty Programs ----

func (h *Handler) CreateLoyaltyProgram(c *fiber.Ctx) error {
	var req crmreq.CreateLoyaltyProgram
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	p, err := h.loyalty.CreateProgram(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, p)
}

func (h *Handler) ListLoyaltyPrograms(c *fiber.Ctx) error {
	programs, err := h.loyalty.ListPrograms(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, programs)
}

// ---- Loyalty Transactions ----

func (h *Handler) EarnPoints(c *fiber.Ctx) error {
	var req crmreq.EarnPoints
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	tx, err := h.loyalty.EarnPoints(c.Context(), c.Params("id"), tenantID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, tx)
}

func (h *Handler) RedeemPoints(c *fiber.Ctx) error {
	var req crmreq.RedeemPoints
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	tx, err := h.loyalty.RedeemPoints(c.Context(), c.Params("id"), tenantID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, tx)
}

func (h *Handler) GetLoyaltyBalance(c *fiber.Ctx) error {
	balance, err := h.loyalty.GetBalance(c.Context(), c.Params("id"), c.Params("program_id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, balance)
}

func (h *Handler) ListLoyaltyTransactions(c *fiber.Ctx) error {
	txs, err := h.loyalty.ListTransactions(c.Context(), c.Params("id"), tenantID(c),
		c.QueryInt("page", 1), c.QueryInt("per_page", 20),
	)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, txs)
}
