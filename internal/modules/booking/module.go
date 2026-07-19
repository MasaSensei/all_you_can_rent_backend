package booking

import (
	"github.com/gofiber/fiber/v2"

	"rentos/internal/bootstrap"
	"rentos/internal/modules/booking/handler"
	"rentos/internal/modules/booking/repository/postgres"
	"rentos/internal/modules/booking/routes"
	"rentos/internal/modules/booking/service"
)

// Module holds the booking module's wired handler.
type Module struct {
	handler *handler.Handler
}

// New builds the booking module. It receives an InventoryChecker (from
// the inventory module) so the booking service can validate asset
// availability without importing the inventory package directly.
// A PassthroughPricer is used until Phase 5 (pricing module) is wired.
func New(c *bootstrap.Container, inventory service.InventoryChecker) *Module {
	bookingRepo := postgres.NewBookingRepository(
		query("create_booking.sql"),
		query("find_booking_by_id.sql"),
		query("list_bookings.sql"),
		query("update_booking_status.sql"),
		query("update_booking_totals.sql"),
	)
	itemRepo := postgres.NewBookingItemRepository(
		query("create_booking_item.sql"),
		query("find_booking_item_by_id.sql"),
		query("list_booking_items.sql"),
		query("update_booking_item_end_date.sql"),
		query("check_overlapping_booking_items.sql"),
	)
	extensionRepo := postgres.NewBookingExtensionRepository(
		query("create_booking_extension.sql"),
		query("list_booking_extensions.sql"),
	)
	returnRepo := postgres.NewBookingReturnRepository(
		query("create_booking_return.sql"),
		query("list_booking_returns.sql"),
	)

	// PassthroughPricer is replaced with the real PricingService in Phase 5.
	pricer := &service.PassthroughPricer{}

	bookingSvc := service.NewBookingService(c.DB, bookingRepo, itemRepo, inventory, pricer)
	itemSvc := service.NewBookingItemService(c.DB, bookingRepo, itemRepo, extensionRepo, returnRepo, inventory, pricer)

	h := handler.New(bookingSvc, itemSvc, c.Validator)
	return &Module{handler: h}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}
