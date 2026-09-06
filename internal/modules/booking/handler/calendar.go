package handler

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/pkg/response"
)

// GET /bookings/calendar?start_date=&end_date=&category_id=&status=
func (h *Handler) Calendar(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(string)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	categoryID := c.Query("category_id")
	status := c.Query("status")

	if startDate == "" || endDate == "" {
		return response.Error(c, response.NewAppError(response.CodeValidation,
			"start_date dan end_date wajib diisi"))
	}

	slots, err := h.bookings.CalendarView(c.Context(), tenantID, startDate, endDate, categoryID, status)
	if err != nil {
		return response.FromError(c, err)
	}
	return response.Success(c, slots)
}
