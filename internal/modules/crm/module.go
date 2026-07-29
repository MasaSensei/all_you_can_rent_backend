package crm

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/crm/handler"
	"rentos-backend/internal/modules/crm/repository/postgres"
	"rentos-backend/internal/modules/crm/routes"
	"rentos-backend/internal/modules/crm/service"
)

// Module holds the CRM module's wired handler.
type Module struct {
	handler *handler.Handler
}

// New builds the CRM module: repositories → services → handler.
func New(c *bootstrap.Container) *Module {
	customerRepo := postgres.NewCustomerRepository(
		query("create_customer.sql"),
		query("find_customer_by_id.sql"),
		query("find_customer_by_email.sql"),
		query("list_customers.sql"),
		query("update_customer.sql"),
		query("delete_customer.sql"),
	)
	addressRepo := postgres.NewCustomerAddressRepository(
		query("add_address.sql"),
		query("list_addresses.sql"),
		query("delete_address.sql"),
		query("unset_default_addresses.sql"),
	)
	membershipRepo := postgres.NewMembershipRepository(
		query("create_membership.sql"),
		query("find_membership_by_id.sql"),
		query("list_memberships_by_customer.sql"),
	)
	loyaltyProgramRepo := postgres.NewLoyaltyProgramRepository(
		query("create_loyalty_program.sql"),
		query("find_loyalty_program_by_id.sql"),
		query("list_loyalty_programs.sql"),
	)
	loyaltyTxRepo := postgres.NewLoyaltyTransactionRepository(
		query("create_loyalty_transaction.sql"),
		query("list_loyalty_transactions.sql"),
		query("sum_loyalty_balance.sql"),
	)

	customerSvc := service.NewCustomerService(c.DB, customerRepo, addressRepo)
	membershipSvc := service.NewMembershipService(c.DB, membershipRepo, customerRepo)
	loyaltySvc := service.NewLoyaltyService(c.DB, loyaltyProgramRepo, loyaltyTxRepo, customerRepo)

	h := handler.New(customerSvc, membershipSvc, loyaltySvc, c.Validator)
	return &Module{handler: h}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}
