package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/crm/handler"
)

// Register mounts all CRM routes under /api/v1.
func Register(router fiber.Router, h *handler.Handler) {
	customers := router.Group("/customers")
	customers.Post("/", h.CreateCustomer)
	customers.Get("/", h.ListCustomers)
	customers.Get("/:id", h.GetCustomer)
	customers.Put("/:id", h.UpdateCustomer)
	customers.Delete("/:id", h.DeleteCustomer)

	// Customer Addresses
	customers.Post("/:id/addresses", h.AddAddress)
	customers.Delete("/:id/addresses/:addr_id", h.DeleteAddress)

	// Customer Memberships
	customers.Post("/:id/memberships", h.CreateMembership)
	customers.Get("/:id/memberships", h.ListMemberships)

	// Customer Loyalty
	customers.Post("/:id/loyalty/earn", h.EarnPoints)
	customers.Post("/:id/loyalty/redeem", h.RedeemPoints)
	customers.Get("/:id/loyalty/balance/:program_id", h.GetLoyaltyBalance)
	customers.Get("/:id/loyalty/transactions", h.ListLoyaltyTransactions)

	// Loyalty Programs (tenant-level)
	router.Post("/loyalty-programs", h.CreateLoyaltyProgram)
	router.Get("/loyalty-programs", h.ListLoyaltyPrograms)
}
