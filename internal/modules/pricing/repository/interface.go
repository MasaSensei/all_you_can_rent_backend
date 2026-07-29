package repository

import (
	"context"
	"errors"

	"rentos-backend/internal/modules/pricing/entity"
	"rentos-backend/pkg/database"
)

var ErrNotFound = errors.New("repository: record not found")

// PricingRuleRepository manages the pricing_rules table.
type PricingRuleRepository interface {
	Create(ctx context.Context, q database.Querier, r *entity.PricingRule) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.PricingRule, error)
	List(ctx context.Context, q database.Querier, tenantID string) ([]entity.PricingRule, error)
	// FindApplicable returns the best matching rule for an asset+category combination.
	// Asset-level rules take priority over category-level rules.
	FindApplicable(ctx context.Context, q database.Querier, tenantID, assetID, categoryID string) (*entity.PricingRule, error)
	Update(ctx context.Context, q database.Querier, r *entity.PricingRule) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// CouponRepository manages the coupons table.
type CouponRepository interface {
	Create(ctx context.Context, q database.Querier, c *entity.Coupon) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Coupon, error)
	FindByCode(ctx context.Context, q database.Querier, code, tenantID string) (*entity.Coupon, error)
	List(ctx context.Context, q database.Querier, tenantID string, limit, offset int) ([]entity.Coupon, error)
	IncrementUsage(ctx context.Context, q database.Querier, id string) error
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}

// PromotionRepository manages the promotions table.
type PromotionRepository interface {
	Create(ctx context.Context, q database.Querier, p *entity.Promotion) error
	FindByID(ctx context.Context, q database.Querier, id, tenantID string) (*entity.Promotion, error)
	ListActive(ctx context.Context, q database.Querier, tenantID string) ([]entity.Promotion, error)
	Delete(ctx context.Context, q database.Querier, id, tenantID string) error
}
