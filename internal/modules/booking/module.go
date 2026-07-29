package booking

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/booking/handler"
	"rentos-backend/internal/modules/booking/repository/postgres"
	"rentos-backend/internal/modules/booking/routes"
	"rentos-backend/internal/modules/booking/service"
)

// Module holds the booking module's wired handler.
type Module struct {
	handler *handler.Handler
}

// New builds the booking module.
//
//   - inventory satisfies service.InventoryChecker
//     (injected from inventory.Module.AssetService())
//   - pricing satisfies service.PricingQuoter
//     (injected from pricing.Module.PricingQuoter())
//
// Both are interface parameters — booking never imports inventory or
// pricing concrete packages.
func New(c *bootstrap.Container, inventory service.InventoryChecker, pricing service.PricingQuoter) *Module {
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

	bookingSvc := service.NewBookingService(c.DB, bookingRepo, itemRepo, inventory, pricing)
	itemSvc := service.NewBookingItemService(c.DB, bookingRepo, itemRepo, extensionRepo, returnRepo, inventory, pricing)

	h := handler.New(bookingSvc, itemSvc, c.Validator)
	return &Module{handler: h}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}
