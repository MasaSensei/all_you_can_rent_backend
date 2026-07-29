package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"rentos-backend/internal/modules/pricing/entity"
	"rentos-backend/internal/modules/pricing/repository"
	"rentos-backend/pkg/database"
)

// ============================================================
// pricingRuleRepository
// ============================================================

type pricingRuleRepository struct {
	qCreate         string
	qFindByID       string
	qList           string
	qFindApplicable string
	qUpdate         string
	qDelete         string
}

func NewPricingRuleRepository(qCreate, qFindByID, qList, qFindApplicable, qUpdate, qDelete string) repository.PricingRuleRepository {
	return &pricingRuleRepository{
		qCreate: qCreate, qFindByID: qFindByID, qList: qList,
		qFindApplicable: qFindApplicable, qUpdate: qUpdate, qDelete: qDelete,
	}
}

func (r *pricingRuleRepository) Create(ctx context.Context, q database.Querier, rule *entity.PricingRule) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		rule.ID, rule.TenantID, rule.CategoryID, rule.AssetID, rule.Name,
		rule.RuleType, rule.Value, rule.DurationUnit, rule.MinDuration,
		rule.MaxDuration, rule.ValidFrom, rule.ValidTo, rule.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("pricingRuleRepository.Create: %w", err)
	}
	return nil
}

func (r *pricingRuleRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.PricingRule, error) {
	var rule entity.PricingRule
	if err := q.GetContext(ctx, &rule, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("pricingRuleRepository.FindByID: %w", err)
	}
	return &rule, nil
}

func (r *pricingRuleRepository) List(ctx context.Context, q database.Querier, tenantID string) ([]entity.PricingRule, error) {
	var out []entity.PricingRule
	if err := q.SelectContext(ctx, &out, r.qList, tenantID); err != nil {
		return nil, fmt.Errorf("pricingRuleRepository.List: %w", err)
	}
	return out, nil
}

func (r *pricingRuleRepository) FindApplicable(ctx context.Context, q database.Querier, tenantID, assetID, categoryID string) (*entity.PricingRule, error) {
	var rule entity.PricingRule
	if err := q.GetContext(ctx, &rule, r.qFindApplicable, tenantID, assetID, categoryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("pricingRuleRepository.FindApplicable: %w", err)
	}
	return &rule, nil
}

func (r *pricingRuleRepository) Update(ctx context.Context, q database.Querier, rule *entity.PricingRule) error {
	res, err := q.ExecContext(ctx, r.qUpdate,
		rule.ID, rule.TenantID, rule.Name, rule.Value, rule.ValidFrom, rule.ValidTo,
	)
	if err != nil {
		return fmt.Errorf("pricingRuleRepository.Update: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *pricingRuleRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("pricingRuleRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// couponRepository
// ============================================================

type couponRepository struct {
	qCreate          string
	qFindByID        string
	qFindByCode      string
	qList            string
	qIncrementUsage  string
	qDelete          string
}

func NewCouponRepository(qCreate, qFindByID, qFindByCode, qList, qIncrementUsage, qDelete string) repository.CouponRepository {
	return &couponRepository{
		qCreate: qCreate, qFindByID: qFindByID, qFindByCode: qFindByCode,
		qList: qList, qIncrementUsage: qIncrementUsage, qDelete: qDelete,
	}
}

func (r *couponRepository) Create(ctx context.Context, q database.Querier, c *entity.Coupon) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		c.ID, c.TenantID, c.Code, c.DiscountType, c.DiscountValue,
		c.MinOrderValue, c.UsageLimit, c.ValidFrom, c.ValidTo, c.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("couponRepository.Create: %w", err)
	}
	return nil
}

func (r *couponRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Coupon, error) {
	var c entity.Coupon
	if err := q.GetContext(ctx, &c, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("couponRepository.FindByID: %w", err)
	}
	return &c, nil
}

func (r *couponRepository) FindByCode(ctx context.Context, q database.Querier, code, tenantID string) (*entity.Coupon, error) {
	var c entity.Coupon
	if err := q.GetContext(ctx, &c, r.qFindByCode, code, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("couponRepository.FindByCode: %w", err)
	}
	return &c, nil
}

func (r *couponRepository) List(ctx context.Context, q database.Querier, tenantID string, limit, offset int) ([]entity.Coupon, error) {
	var out []entity.Coupon
	if err := q.SelectContext(ctx, &out, r.qList, tenantID, limit, offset); err != nil {
		return nil, fmt.Errorf("couponRepository.List: %w", err)
	}
	return out, nil
}

func (r *couponRepository) IncrementUsage(ctx context.Context, q database.Querier, id string) error {
	_, err := q.ExecContext(ctx, r.qIncrementUsage, id)
	if err != nil {
		return fmt.Errorf("couponRepository.IncrementUsage: %w", err)
	}
	return nil
}

func (r *couponRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("couponRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// ============================================================
// promotionRepository
// ============================================================

type promotionRepository struct {
	qCreate     string
	qFindByID   string
	qListActive string
	qDelete     string
}

func NewPromotionRepository(qCreate, qFindByID, qListActive, qDelete string) repository.PromotionRepository {
	return &promotionRepository{
		qCreate: qCreate, qFindByID: qFindByID,
		qListActive: qListActive, qDelete: qDelete,
	}
}

func (r *promotionRepository) Create(ctx context.Context, q database.Querier, p *entity.Promotion) error {
	_, err := q.ExecContext(ctx, r.qCreate,
		p.ID, p.TenantID, p.Name, p.Description, p.PromotionType,
		p.Value, p.StartDate, p.EndDate, p.CreatedBy,
	)
	if err != nil {
		return fmt.Errorf("promotionRepository.Create: %w", err)
	}
	return nil
}

func (r *promotionRepository) FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Promotion, error) {
	var p entity.Promotion
	if err := q.GetContext(ctx, &p, r.qFindByID, id, tenantID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("promotionRepository.FindByID: %w", err)
	}
	return &p, nil
}

func (r *promotionRepository) ListActive(ctx context.Context, q database.Querier, tenantID string) ([]entity.Promotion, error) {
	var out []entity.Promotion
	if err := q.SelectContext(ctx, &out, r.qListActive, tenantID); err != nil {
		return nil, fmt.Errorf("promotionRepository.ListActive: %w", err)
	}
	return out, nil
}

func (r *promotionRepository) Delete(ctx context.Context, q database.Querier, id, tenantID string) error {
	res, err := q.ExecContext(ctx, r.qDelete, id, tenantID)
	if err != nil {
		return fmt.Errorf("promotionRepository.Delete: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}
