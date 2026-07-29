package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/booking/handler"
)

// Register mounts all booking routes under /api/v1.
func Register(router fiber.Router, h *handler.Handler) {
	bookings := router.Group("/bookings")
	bookings.Post("/", h.CreateBooking)
	bookings.Get("/", h.ListBookings)
	bookings.Get("/:id", h.GetBooking)
	bookings.Post("/:id/confirm", h.ConfirmBooking)
	bookings.Post("/:id/cancel", h.CancelBooking)

	items := router.Group("/booking-items")
	items.Post("/:item_id/extend", h.ExtendBookingItem)
	items.Post("/:item_id/return", h.ReturnBookingItem)
}
