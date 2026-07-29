package service

import (
	"context"

	"rentos-backend/internal/modules/crm/dto/request"
	"rentos-backend/internal/modules/crm/dto/response"
	"rentos-backend/internal/modules/crm/entity"
)

// CustomerService manages customer lifecycle.
type CustomerService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateCustomer) (*response.Customer, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Customer, error)
	List(ctx context.Context, tenantID string, filter request.ListCustomersFilter) ([]response.Customer, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateCustomer) (*response.Customer, error)
	Delete(ctx context.Context, id, tenantID string) error

	AddAddress(ctx context.Context, customerID, tenantID, actorID string, req request.AddAddress) (*entity.CustomerAddress, error)
	DeleteAddress(ctx context.Context, addressID, customerID, tenantID, actorID string) error
}

// MembershipService manages customer memberships.
type MembershipService interface {
	Create(ctx context.Context, customerID, tenantID, actorID string, req request.CreateMembership) (*response.Membership, error)
	GetByID(ctx context.Context, id, tenantID string) (*response.Membership, error)
	ListByCustomer(ctx context.Context, customerID, tenantID string) ([]response.Membership, error)
}

// LoyaltyService manages loyalty programs and transactions.
type LoyaltyService interface {
	CreateProgram(ctx context.Context, tenantID, actorID string, req request.CreateLoyaltyProgram) (*response.LoyaltyProgram, error)
	ListPrograms(ctx context.Context, tenantID string) ([]response.LoyaltyProgram, error)

	EarnPoints(ctx context.Context, customerID, tenantID string, req request.EarnPoints) (*response.LoyaltyTransaction, error)
	RedeemPoints(ctx context.Context, customerID, tenantID string, req request.RedeemPoints) (*response.LoyaltyTransaction, error)
	GetBalance(ctx context.Context, customerID, loyaltyProgramID, tenantID string) (*response.LoyaltyBalance, error)
	ListTransactions(ctx context.Context, customerID, tenantID string, page, perPage int) ([]response.LoyaltyTransaction, error)
}
