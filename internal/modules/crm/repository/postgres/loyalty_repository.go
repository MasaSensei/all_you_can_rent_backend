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
// loyaltyProgramRepository
// ============================================================

type loyaltyProgramRepository struct {
	qCreate   string
	qFindByID string
	qList     string
}

func NewLoyaltyProgramRepository(qCreate, qFindByID, qList string) repository.LoyaltyProgramRepository {
	return &loyaltyProgramRepository{qCreate: qCreate, qFindByID: qFindByID, qList: qList}
}

func (r *loyaltyProgramRepository) Create(ctx context.Context, q database.Querier, p *entity.LoyaltyProgram) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		p.ID, p.TenantID, p.Name, p.Description,
		p.PointsPerCurrency, p.RedemptionRate, p.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("loyaltyProgramRepository.Create: %w", err)
	}
	return nil
}

func (r *loyaltyProgramRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.LoyaltyProgram, error) {
	var p entity.LoyaltyProgram
	if err := q.GetContext(ctx, &p, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("loyaltyProgramRepository.FindByID: %w", err)
	}
	return &p, nil
}

func (r *loyaltyProgramRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.LoyaltyProgram, error) {
	var out []entity.LoyaltyProgram
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("loyaltyProgramRepository.List: %w", err)
	}
	return out, nil
}

// ============================================================
// loyaltyTransactionRepository
// ============================================================

type loyaltyTransactionRepository struct {
	qCreate        string
	qListByCustomer string
	qSumBalance    string
}

func NewLoyaltyTransactionRepository(qCreate, qListByCustomer, qSumBalance string) repository.LoyaltyTransactionRepository {
	return &loyaltyTransactionRepository{
		qCreate: qCreate, qListByCustomer: qListByCustomer, qSumBalance: qSumBalance,
	}
}

func (r *loyaltyTransactionRepository) Create(ctx context.Context, q database.Querier, t *entity.LoyaltyTransaction) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		t.ID, t.TenantID, t.LoyaltyProgramID, t.CustomerID, t.BookingID,
		t.Points, t.TransactionType, t.Description,
	)
	if err != nil {
		return fmt.Errorf("loyaltyTransactionRepository.Create: %w", err)
	}
	return nil
}

func (r *loyaltyTransactionRepository) ListByCustomer(ctx context.Context, q database.Querier, customerID, tenantID string, limit, offset int) ([]entity.LoyaltyTransaction, error) {
	var out []entity.LoyaltyTransaction
	if err := q.SelectContext(ctx, &out, r.qListByCustomer, customerID, tenantID, limit, offset); err != nil {
		return nil, fmt.Errorf("loyaltyTransactionRepository.ListByCustomer: %w", err)
	}
	return out, nil
}

func (r *loyaltyTransactionRepository) SumBalance(ctx context.Context, q database.Querier, customerID, loyaltyProgramID string) (int, error) {
	var balance int
	if err := q.GetContext(ctx, &balance, r.qSumBalance, customerID, loyaltyProgramID); err != nil {
		return 0, fmt.Errorf("loyaltyTransactionRepository.SumBalance: %w", err)
	}
	return balance, nil
}
