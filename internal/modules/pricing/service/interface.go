package service

import (
	"context"
	"time"

	"rentos-backend/internal/modules/pricing/dto/request"
	"rentos-backend/internal/modules/pricing/dto/response"
	"rentos-backend/internal/modules/pricing/entity"
)

// PricingRuleService manages pricing rules.
type PricingRuleService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreatePricingRule) (*entity.PricingRule, error)
	GetByID(ctx context.Context, id, tenantID string) (*entity.PricingRule, error)
	List(ctx context.Context, tenantID string) ([]entity.PricingRule, error)
	Update(ctx context.Context, id, tenantID, actorID string, req request.UpdatePricingRule) (*entity.PricingRule, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// CouponService manages coupons.
type CouponService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreateCoupon) (*entity.Coupon, error)
	GetByID(ctx context.Context, id, tenantID string) (*entity.Coupon, error)
	List(ctx context.Context, tenantID string, page, perPage int) ([]entity.Coupon, error)
	Delete(ctx context.Context, id, tenantID string) error
	Validate(ctx context.Context, tenantID string, req request.ValidateCoupon) (*response.CouponValidation, error)
}

// PromotionService manages promotions.
type PromotionService interface {
	Create(ctx context.Context, tenantID, actorID string, req request.CreatePromotion) (*entity.Promotion, error)
	GetByID(ctx context.Context, id, tenantID string) (*entity.Promotion, error)
	ListActive(ctx context.Context, tenantID string) ([]entity.Promotion, error)
	Delete(ctx context.Context, id, tenantID string) error
}

// PricingQuoter is implemented by PricingService and satisfies the
// booking.PricingQuoter interface, replacing PassthroughPricer in
// Phase 5. Defined here so the pricing package is self-contained and
// the booking module only imports its own service/interface.go.
type PricingQuoter interface {
	// QuoteItem returns the unit price for one asset over a given date range.
	// It resolves the best matching pricing rule (asset > category > default).
	QuoteItem(ctx context.Context, tenantID, assetID string, start, end time.Time, qty int) (float64, error)

	// ValidateCoupon checks a coupon code, returns discount amount + coupon ID.
	// Returns (0, nil, nil) when code is empty — safe to call unconditionally.
	ValidateCoupon(ctx context.Context, tenantID, code string, subtotal float64) (discount float64, couponID *string, err error)
}
