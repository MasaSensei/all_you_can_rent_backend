package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos-backend/internal/modules/finance/entity"
	"rentos-backend/internal/modules/finance/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// refundRepository
// ============================================================

type refundRepository struct {
	qCreate        string
	qFindByID      string
	qSumByPayment  string
}

func NewRefundRepository(qCreate, qFindByID, qSumByPayment string) repository.RefundRepository {
	return &refundRepository{
		qCreate: qCreate, qFindByID: qFindByID, qSumByPayment: qSumByPayment,
	}
}

func (r *refundRepository) Create(ctx context.Context, q database.Querier, ref *entity.Refund) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		ref.ID, ref.TenantID, ref.PaymentID, ref.Amount, ref.Reason, ref.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("refundRepository.Create: %w", err)
	}
	return nil
}

func (r *refundRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Refund, error) {
	var ref entity.Refund
	if err := q.GetContext(ctx, &ref, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("refundRepository.FindByID: %w", err)
	}
	return &ref, nil
}

func (r *refundRepository) SumByPayment(ctx context.Context, q database.Querier, paymentID string) (float64, error) {
	var total float64
	if err := q.GetContext(ctx, &total, r.qSumByPayment, paymentID); err != nil {
		return 0, fmt.Errorf("refundRepository.SumByPayment: %w", err)
	}
	return total, nil
}

// ============================================================
// taxRepository
// ============================================================

type taxRepository struct {
	qCreate      string
	qFindDefault string
	qList        string
}

func NewTaxRepository(qCreate, qFindDefault, qList string) repository.TaxRepository {
	return &taxRepository{qCreate: qCreate, qFindDefault: qFindDefault, qList: qList}
}

func (r *taxRepository) Create(ctx context.Context, q database.Querier, t *entity.Tax) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		t.ID, t.TenantID, t.Name, t.Rate, t.TaxType, t.IsDefault, t.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("taxRepository.Create: %w", err)
	}
	return nil
}

func (r *taxRepository) FindDefault(ctx context.Context, q database.Querier, tenantID string) (*entity.Tax, error) {
	var t entity.Tax
	if err := q.GetContext(ctx, &t, r.qFindDefault, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("taxRepository.FindDefault: %w", err)
	}
	return &t, nil
}

func (r *taxRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.Tax, error) {
	var out []entity.Tax
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("taxRepository.List: %w", err)
	}
	return out, nil
}

// ============================================================
// expenseRepository
// ============================================================

type expenseRepository struct {
	qCreate string
	qList   string
}

func NewExpenseRepository(qCreate, qList string) repository.ExpenseRepository {
	return &expenseRepository{qCreate: qCreate, qList: qList}
}

func (r *expenseRepository) Create(ctx context.Context, q database.Querier, e *entity.Expense) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		e.ID, e.TenantID, e.AssetID, e.Category, e.Amount, e.ExpenseDate,
		e.Description, e.Vendor, e.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("expenseRepository.Create: %w", err)
	}
	return nil
}

func (r *expenseRepository) List(ctx context.Context, q database.Querier, tenantID string, assetID, category *string, limit, offset int) ([]entity.Expense, error) {
	var out []entity.Expense
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, assetID, category, limit, offset); err != nil {
		return nil, fmt.Errorf("expenseRepository.List: %w", err)
	}
	return out, nil
}
