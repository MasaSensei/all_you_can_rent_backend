package handler

import (
	"github.com/gofiber/fiber/v2"

	bookingreq "rentos-backend/internal/modules/booking/dto/request"
	"rentos-backend/internal/modules/booking/service"
	apiresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

const (
	ctxKeyTenantID = "tenant_id"
	ctxKeyUserID   = "user_id"
)

// Handler groups the booking module's HTTP handlers.
type Handler struct {
	bookings service.BookingService
	items    service.BookingItemService
	validate *validator.Validate
}

func New(
	bookings service.BookingService,
	items service.BookingItemService,
	v *validator.Validate,
) *Handler {
	return &Handler{bookings: bookings, items: items, validate: v}
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

// ---- Bookings ----

func (h *Handler) CreateBooking(c *fiber.Ctx) error {
	var req bookingreq.CreateBooking
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	b, err := h.bookings.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, b)
}

func (h *Handler) GetBooking(c *fiber.Ctx) error {
	b, err := h.bookings.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, b)
}

func (h *Handler) ListBookings(c *fiber.Ctx) error {
	filter := bookingreq.ListBookingsFilter{
		Page:    c.QueryInt("page", 1),
		PerPage: c.QueryInt("per_page", 20),
	}
	if v := c.Query("customer_id"); v != "" {
		filter.CustomerID = &v
	}
	if v := c.Query("booking_status"); v != "" {
		filter.BookingStatus = &v
	}

	bookings, err := h.bookings.List(c.Context(), tenantID(c), filter)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, bookings, fiber.Map{
		"page": filter.Page, "per_page": filter.PerPage,
	})
}

func (h *Handler) ConfirmBooking(c *fiber.Ctx) error {
	b, err := h.bookings.Confirm(c.Context(), c.Params("id"), tenantID(c), userID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, b)
}

func (h *Handler) CancelBooking(c *fiber.Ctx) error {
	var req bookingreq.CancelBooking
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}

	b, err := h.bookings.Cancel(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, b)
}

// ---- Booking Items ----

func (h *Handler) ExtendBookingItem(c *fiber.Ctx) error {
	var req bookingreq.ExtendBookingItem
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	ext, err := h.items.Extend(c.Context(), c.Params("item_id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, ext)
}

func (h *Handler) ReturnBookingItem(c *fiber.Ctx) error {
	var req bookingreq.ReturnBookingItem
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}

	ret, err := h.items.Return(c.Context(), c.Params("item_id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, ret)
}
