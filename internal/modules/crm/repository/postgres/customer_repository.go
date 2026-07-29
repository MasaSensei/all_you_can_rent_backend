package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos-backend/internal/modules/crm/entity"
	"rentos-backend/internal/modules/crm/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// customerRepository
// ============================================================

type customerRepository struct {
	qCreate        string
	qFindByID      string
	qFindByEmail   string
	qList          string
	qUpdate        string
	qDelete        string
}

func NewCustomerRepository(qCreate, qFindByID, qFindByEmail, qList, qUpdate, qDelete string) repository.CustomerRepository {
	return &customerRepository{
		qCreate: qCreate, qFindByID: qFindByID, qFindByEmail: qFindByEmail,
		qList: qList, qUpdate: qUpdate, qDelete: qDelete,
	}
}

func (r *customerRepository) Create(ctx context.Context, q database.Querier, c *entity.Customer) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		c.ID, c.TenantID, c.FirstName, c.LastName, c.Email, c.Phone, c.CompanyName,
		c.DateOfBirth, c.IDDocumentType, c.IDDocumentNumber,
		c.CustomerType, c.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("customerRepository.Create: %w", err)
	}
	return nil
}

func (r *customerRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Customer, error) {
	var c entity.Customer
	if err := q.GetContext(ctx, &c, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("customerRepository.FindByID: %w", err)
	}
	return &c, nil
}

func (r *customerRepository) FindByEmail(ctx context.Context, q database.Querier, email, tenantID string) (*entity.Customer, error) {
	var c entity.Customer
	if err := q.GetContext(ctx, &c, r.qFindByEmail, email, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("customerRepository.FindByEmail: %w", err)
	}
	return &c, nil
}

func (r *customerRepository) List(ctx context.Context, q database.Querier, tenantID string, customerType, search *string, limit, offset int) ([]entity.Customer, error) {
	var out []entity.Customer
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, customerType, search, limit, offset); err != nil {
		return nil, fmt.Errorf("customerRepository.List: %w", err)
	}
	return out, nil
}

func (r *customerRepository) Update(ctx context.Context, q database.Querier, c *entity.Customer) error {
	res, err := q.ExecContext(ctx, r.qUpdate,
		c.ID, c.TenantID, c.FirstName, c.LastName, c.Phone, c.CompanyName,
		c.DateOfBirth, c.IDDocumentType, c.IDDocumentNumber,
	)
	if err != nil {
		return fmt.Errorf("customerRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *customerRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("customerRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// customerAddressRepository
// ============================================================

type customerAddressRepository struct {
	qAdd              string
	qList             string
	qDelete           string
	qUnsetAllDefaults string
}

func NewCustomerAddressRepository(qAdd, qList, qDelete, qUnsetAllDefaults string) repository.CustomerAddressRepository {
	return &customerAddressRepository{
		qAdd: qAdd, qList: qList, qDelete: qDelete, qUnsetAllDefaults: qUnsetAllDefaults,
	}
}

func (r *customerAddressRepository) Add(ctx context.Context, q database.Querier, a *entity.CustomerAddress) error {
	_, err := q.ExecContext(ctx, r.qAdd,
		a.ID, a.TenantID, a.CustomerID, a.AddressType, a.Line1, a.Line2,
		a.City, a.State, a.PostalCode, a.Country, a.IsDefault, a.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("customerAddressRepository.Add: %w", err)
	}
	return nil
}

func (r *customerAddressRepository) ListByCustomer(ctx context.Context, q database.Querier, customerID, tenantID string) ([]entity.CustomerAddress, error) {
	var out []entity.CustomerAddress
	if err := q.SelectContext(ctx, &out, r.qList, customerID, tenantID); err != nil {
		return nil, fmt.Errorf("customerAddressRepository.ListByCustomer: %w", err)
	}
	return out, nil
}

func (r *customerAddressRepository) Delete(ctx context.Context, q database.Querier, id, customerID, actorID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, customerID, actorID)
	if err != nil {
		return fmt.Errorf("customerAddressRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *customerAddressRepository) UnsetAllDefaults(ctx context.Context, q database.Querier, customerID, tenantID string) error {
	_, err := q.ExecContext(ctx, r.qUnsetAllDefaults, customerID, tenantID)
	if err != nil {
		return fmt.Errorf("customerAddressRepository.UnsetAllDefaults: %w", err)
	}
	return nil
}

// ============================================================
// membershipRepository
// ============================================================

type membershipRepository struct {
	qCreate        string
	qFindByID      string
	qListByCustomer string
}

func NewMembershipRepository(qCreate, qFindByID, qListByCustomer string) repository.MembershipRepository {
	return &membershipRepository{qCreate: qCreate, qFindByID: qFindByID, qListByCustomer: qListByCustomer}
}

func (r *membershipRepository) Create(ctx context.Context, q database.Querier, m *entity.Membership) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		m.ID, m.TenantID, m.CustomerID, m.PlanName, m.Tier,
		m.StartDate, m.EndDate, m.Fee, m.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("membershipRepository.Create: %w", err)
	}
	return nil
}

func (r *membershipRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Membership, error) {
	var m entity.Membership
	if err := q.GetContext(ctx, &m, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("membershipRepository.FindByID: %w", err)
	}
	return &m, nil
}

func (r *membershipRepository) ListByCustomer(ctx context.Context, q database.Querier, customerID, tenantID string) ([]entity.Membership, error) {
	var out []entity.Membership
	if err := q.SelectContext(ctx, &out, r.qListByCustomer, customerID, tenantID); err != nil {
		return nil, fmt.Errorf("membershipRepository.ListByCustomer: %w", err)
	}
	return out, nil
}
