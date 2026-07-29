package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/crm/dto/request"
	"rentos-backend/internal/modules/crm/dto/response"
	"rentos-backend/internal/modules/crm/entity"
	"rentos-backend/internal/modules/crm/repository"
	pkgresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/transaction"
)

// ============================================================
// customerService
// ============================================================

type customerService struct {
	db        *sqlx.DB
	customers repository.CustomerRepository
	addresses repository.CustomerAddressRepository
}

func NewCustomerService(
	db *sqlx.DB,
	customers repository.CustomerRepository,
	addresses repository.CustomerAddressRepository,
) CustomerService {
	return &customerService{db: db, customers: customers, addresses: addresses}
}

func (s *customerService) Create(ctx context.Context, tenantID, actorID string, req request.CreateCustomer) (*response.Customer, error) {
	if _, err := s.customers.FindByEmail(ctx, s.db, req.Email, tenantID); !errors.Is(err, repository.ErrNotFound) {
		if err == nil {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "a customer with this email already exists")
		}
		return nil, err
	}

	c := &entity.Customer{
		ID:               uuid.NewString(),
		TenantID:         tenantID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            req.Email,
		Phone:            req.Phone,
		CompanyName:      req.CompanyName,
		DateOfBirth:      req.DateOfBirth,
		IDDocumentType:   req.IDDocumentType,
		IDDocumentNumber: req.IDDocumentNumber,
		CustomerType:     req.CustomerType,
		CreatedBy:        &actorID,
	}
	if err := s.customers.Create(ctx, s.db, c); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, c.ID, tenantID)
}

func (s *customerService) GetByID(ctx context.Context, id, tenantID string) (*response.Customer, error) {
	c, err := s.customers.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "customer not found")
		}
		return nil, err
	}
	return s.hydrate(ctx, c)
}

func (s *customerService) List(ctx context.Context, tenantID string, filter request.ListCustomersFilter) ([]response.Customer, error) {
	perPage, page := normPage(filter.PerPage, filter.Page)
	customers, err := s.customers.List(ctx, s.db, tenantID, filter.CustomerType, filter.Search, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.Customer, 0, len(customers))
	for _, c := range customers {
		r, err := s.hydrate(ctx, &c)
		if err != nil {
			return nil, err
		}
		out = append(out, *r)
	}
	return out, nil
}

func (s *customerService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdateCustomer) (*response.Customer, error) {
	c := &entity.Customer{
		ID:               id,
		TenantID:         tenantID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Phone:            req.Phone,
		CompanyName:      req.CompanyName,
		DateOfBirth:      req.DateOfBirth,
		IDDocumentType:   req.IDDocumentType,
		IDDocumentNumber: req.IDDocumentNumber,
		UpdatedBy:        &actorID,
	}
	if err := s.customers.Update(ctx, s.db, c); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "customer not found")
		}
		return nil, err
	}
	return s.GetByID(ctx, id, tenantID)
}

func (s *customerService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.customers.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "customer not found")
		}
		return err
	}
	return nil
}

func (s *customerService) AddAddress(ctx context.Context, customerID, tenantID, actorID string, req request.AddAddress) (*entity.CustomerAddress, error) {
	// Validate customer exists.
	if _, err := s.customers.FindByID(ctx, s.db, customerID, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "customer not found")
		}
		return nil, err
	}

	a := &entity.CustomerAddress{
		ID:          uuid.NewString(),
		TenantID:    tenantID,
		CustomerID:  customerID,
		AddressType: req.AddressType,
		Line1:       req.Line1,
		Line2:       req.Line2,
		City:        req.City,
		State:       req.State,
		PostalCode:  req.PostalCode,
		Country:     req.Country,
		IsDefault:   req.IsDefault,
		CreatedBy:   &actorID,
	}

	if err := transaction.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		// Unset existing defaults before setting a new one.
		if req.IsDefault {
			if err := s.addresses.UnsetAllDefaults(ctx, tx, customerID, tenantID); err != nil {
				return err
			}
		}
		return s.addresses.Add(ctx, tx, a)
	}); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *customerService) DeleteAddress(ctx context.Context, addressID, customerID, tenantID, actorID string) error {
	if err := s.addresses.Delete(ctx, s.db, addressID, customerID, actorID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "address not found")
		}
		return err
	}
	return nil
}

func (s *customerService) hydrate(ctx context.Context, c *entity.Customer) (*response.Customer, error) {
	addrs, _ := s.addresses.ListByCustomer(ctx, s.db, c.ID, c.TenantID)
	addrResp := make([]response.CustomerAddress, 0, len(addrs))
	for _, a := range addrs {
		addrResp = append(addrResp, response.CustomerAddress{
			ID: a.ID, AddressType: a.AddressType, Line1: a.Line1, Line2: a.Line2,
			City: a.City, State: a.State, PostalCode: a.PostalCode,
			Country: a.Country, IsDefault: a.IsDefault,
		})
	}
	return &response.Customer{
		ID: c.ID, TenantID: c.TenantID, FirstName: c.FirstName, LastName: c.LastName,
		Email: c.Email, Phone: c.Phone, CompanyName: c.CompanyName,
		DateOfBirth: c.DateOfBirth, IDDocumentType: c.IDDocumentType,
		IDDocumentNumber: c.IDDocumentNumber, CustomerType: c.CustomerType,
		Status: c.Status, Addresses: addrResp,
		CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}, nil
}

// ============================================================
// membershipService
// ============================================================

type membershipService struct {
	db          *sqlx.DB
	memberships repository.MembershipRepository
	customers   repository.CustomerRepository
}

func NewMembershipService(db *sqlx.DB, memberships repository.MembershipRepository, customers repository.CustomerRepository) MembershipService {
	return &membershipService{db: db, memberships: memberships, customers: customers}
}

func (s *membershipService) Create(ctx context.Context, customerID, tenantID, actorID string, req request.CreateMembership) (*response.Membership, error) {
	if _, err := s.customers.FindByID(ctx, s.db, customerID, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "customer not found")
		}
		return nil, err
	}

	m := &entity.Membership{
		ID:         uuid.NewString(),
		TenantID:   tenantID,
		CustomerID: customerID,
		PlanName:   req.PlanName,
		Tier:       req.Tier,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Fee:        req.Fee,
		CreatedBy:  &actorID,
	}
	if err := s.memberships.Create(ctx, s.db, m); err != nil {
		return nil, err
	}
	return toMembershipResponse(m), nil
}

func (s *membershipService) GetByID(ctx context.Context, id, tenantID string) (*response.Membership, error) {
	m, err := s.memberships.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "membership not found")
		}
		return nil, err
	}
	return toMembershipResponse(m), nil
}

func (s *membershipService) ListByCustomer(ctx context.Context, customerID, tenantID string) ([]response.Membership, error) {
	members, err := s.memberships.ListByCustomer(ctx, s.db, customerID, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.Membership, 0, len(members))
	for _, m := range members {
		out = append(out, *toMembershipResponse(&m))
	}
	return out, nil
}

// ============================================================
// loyaltyService
// ============================================================

type loyaltyService struct {
	db           *sqlx.DB
	programs     repository.LoyaltyProgramRepository
	transactions repository.LoyaltyTransactionRepository
	customers    repository.CustomerRepository
}

func NewLoyaltyService(
	db *sqlx.DB,
	programs repository.LoyaltyProgramRepository,
	transactions repository.LoyaltyTransactionRepository,
	customers repository.CustomerRepository,
) LoyaltyService {
	return &loyaltyService{db: db, programs: programs, transactions: transactions, customers: customers}
}

func (s *loyaltyService) CreateProgram(ctx context.Context, tenantID, actorID string, req request.CreateLoyaltyProgram) (*response.LoyaltyProgram, error) {
	p := &entity.LoyaltyProgram{
		ID:                uuid.NewString(),
		TenantID:          tenantID,
		Name:              req.Name,
		Description:       req.Description,
		PointsPerCurrency: req.PointsPerCurrency,
		RedemptionRate:    req.RedemptionRate,
		CreatedBy:         &actorID,
	}
	if err := s.programs.Create(ctx, s.db, p); err != nil {
		return nil, err
	}
	return toLoyaltyProgramResponse(p), nil
}

func (s *loyaltyService) ListPrograms(ctx context.Context, tenantID string) ([]response.LoyaltyProgram, error) {
	programs, err := s.programs.List(ctx, s.db, tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]response.LoyaltyProgram, 0, len(programs))
	for _, p := range programs {
		out = append(out, *toLoyaltyProgramResponse(&p))
	}
	return out, nil
}

func (s *loyaltyService) EarnPoints(ctx context.Context, customerID, tenantID string, req request.EarnPoints) (*response.LoyaltyTransaction, error) {
	if _, err := s.customers.FindByID(ctx, s.db, customerID, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "customer not found")
		}
		return nil, err
	}
	if _, err := s.programs.FindByID(ctx, s.db, req.LoyaltyProgramID, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "loyalty program not found")
		}
		return nil, err
	}

	desc := req.Description
	t := &entity.LoyaltyTransaction{
		ID:               uuid.NewString(),
		TenantID:         tenantID,
		LoyaltyProgramID: req.LoyaltyProgramID,
		CustomerID:       customerID,
		BookingID:        req.BookingID,
		Points:           req.Points,
		TransactionType:  entity.LoyaltyTxTypeEarn,
		Description:      desc,
	}
	if err := s.transactions.Create(ctx, s.db, t); err != nil {
		return nil, err
	}
	return toLoyaltyTxResponse(t), nil
}

func (s *loyaltyService) RedeemPoints(ctx context.Context, customerID, tenantID string, req request.RedeemPoints) (*response.LoyaltyTransaction, error) {
	if _, err := s.programs.FindByID(ctx, s.db, req.LoyaltyProgramID, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "loyalty program not found")
		}
		return nil, err
	}

	// Guard: enough points.
	balance, err := s.transactions.SumBalance(ctx, s.db, customerID, req.LoyaltyProgramID)
	if err != nil {
		return nil, err
	}
	if req.Points > balance {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "insufficient loyalty points")
	}

	t := &entity.LoyaltyTransaction{
		ID:               uuid.NewString(),
		TenantID:         tenantID,
		LoyaltyProgramID: req.LoyaltyProgramID,
		CustomerID:       customerID,
		Points:           req.Points,
		TransactionType:  entity.LoyaltyTxTypeRedeem,
		Description:      req.Description,
	}
	if err := s.transactions.Create(ctx, s.db, t); err != nil {
		return nil, err
	}
	return toLoyaltyTxResponse(t), nil
}

func (s *loyaltyService) GetBalance(ctx context.Context, customerID, loyaltyProgramID, tenantID string) (*response.LoyaltyBalance, error) {
	balance, err := s.transactions.SumBalance(ctx, s.db, customerID, loyaltyProgramID)
	if err != nil {
		return nil, err
	}
	return &response.LoyaltyBalance{CustomerID: customerID, Balance: balance}, nil
}

func (s *loyaltyService) ListTransactions(ctx context.Context, customerID, tenantID string, page, perPage int) ([]response.LoyaltyTransaction, error) {
	perPage, page = normPage(perPage, page)
	txs, err := s.transactions.ListByCustomer(ctx, s.db, customerID, tenantID, perPage, (page-1)*perPage)
	if err != nil {
		return nil, err
	}
	out := make([]response.LoyaltyTransaction, 0, len(txs))
	for _, t := range txs {
		out = append(out, *toLoyaltyTxResponse(&t))
	}
	return out, nil
}

// ============================================================
// helpers
// ============================================================

func toMembershipResponse(m *entity.Membership) *response.Membership {
	return &response.Membership{
		ID: m.ID, CustomerID: m.CustomerID, PlanName: m.PlanName, Tier: m.Tier,
		StartDate: m.StartDate, EndDate: m.EndDate, Fee: m.Fee,
		MembershipStatus: m.MembershipStatus,
		CreatedAt:        m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func toLoyaltyProgramResponse(p *entity.LoyaltyProgram) *response.LoyaltyProgram {
	return &response.LoyaltyProgram{
		ID: p.ID, TenantID: p.TenantID, Name: p.Name,
		Description: p.Description, PointsPerCurrency: p.PointsPerCurrency,
		RedemptionRate: p.RedemptionRate, Status: p.Status,
	}
}

func toLoyaltyTxResponse(t *entity.LoyaltyTransaction) *response.LoyaltyTransaction {
	return &response.LoyaltyTransaction{
		ID: t.ID, CustomerID: t.CustomerID, BookingID: t.BookingID,
		Points: t.Points, TransactionType: t.TransactionType,
		Description: t.Description, CreatedAt: t.CreatedAt,
	}
}

func normPage(perPage, page int) (int, int) {
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	if page <= 0 {
		page = 1
	}
	return perPage, page
}
