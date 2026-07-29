package repository

import (
	"context"
	"errors"

	"rentos-backend/internal/modules/crm/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

// CustomerRepository manages the customers table.
type CustomerRepository interface {
	Create(ctx context.Context, q database.Querier, c *entity.Customer) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Customer, error)
	FindByEmail(ctx context.Context, q database.Querier, email, tenantID string) (*entity.Customer, error)
	List(ctx context.Context, q database.Querier, tenantID string, customerType, search *string, limit, offset int) ([]entity.Customer, error)
	Update(ctx context.Context, q database.Querier, c *entity.Customer) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// CustomerAddressRepository manages the customer_addresses table.
type CustomerAddressRepository interface {
	Add(ctx context.Context, q database.Querier, a *entity.CustomerAddress) error
	ListByCustomer(ctx context.Context, q database.Querier, customerID, tenantID string) ([]entity.CustomerAddress, error)
	Delete(ctx context.Context, q database.Querier, id, customerID, actorID string) error
	UnsetAllDefaults(ctx context.Context, q database.Querier, customerID, tenantID string) error
}

// MembershipRepository manages the memberships table.
type MembershipRepository interface {
	Create(ctx context.Context, q database.Querier, m *entity.Membership) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Membership, error)
	ListByCustomer(ctx context.Context, q database.Querier, customerID, tenantID string) ([]entity.Membership, error)
}

// LoyaltyProgramRepository manages the loyalty_programs table.
type LoyaltyProgramRepository interface {
	Create(ctx context.Context, q database.Querier, p *entity.LoyaltyProgram) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.LoyaltyProgram, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.LoyaltyProgram, error)
}

// LoyaltyTransactionRepository manages the loyalty_transactions table.
type LoyaltyTransactionRepository interface {
	Create(ctx context.Context, q database.Querier, t *entity.LoyaltyTransaction) error
	ListByCustomer(ctx context.Context, q database.Querier, customerID, tenantID string, limit, offset int) ([]entity.LoyaltyTransaction, error)
	// SumBalance returns the net point balance for a customer in a program.
	SumBalance(ctx context.Context, q database.Querier, customerID, loyaltyProgramID string) (int, error)
}
