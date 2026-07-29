package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"rentos-backend/internal/modules/pricing/dto/request"
	"rentos-backend/internal/modules/pricing/dto/response"
	"rentos-backend/internal/modules/pricing/entity"
	"rentos-backend/internal/modules/pricing/repository"
	pkgresponse "rentos-backend/pkg/response"
)

// ============================================================
// pricingRuleService
// ============================================================

type pricingRuleService struct {
	db   *sqlx.DB
	repo repository.PricingRuleRepository
}

func NewPricingRuleService(db *sqlx.DB, repo repository.PricingRuleRepository) PricingRuleService {
	return &pricingRuleService{db: db, repo: repo}
}

func (s *pricingRuleService) Create(ctx context.Context, tenantID, actorID string, req request.CreatePricingRule) (*entity.PricingRule, error) {
	if req.AssetID == nil && req.CategoryID == nil {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeValidation, "pricing rule must target either an asset or a category")
	}

	rule := &entity.PricingRule{
		ID:           uuid.NewString(),
		TenantID:     tenantID,
		CategoryID:   req.CategoryID,
		AssetID:      req.AssetID,
		Name:         req.Name,
		RuleType:     req.RuleType,
		Value:        req.Value,
		DurationUnit: req.DurationUnit,
		MinDuration:  req.MinDuration,
		MaxDuration:  req.MaxDuration,
		ValidFrom:    req.ValidFrom,
		ValidTo:      req.ValidTo,
		CreatedBy:    &actorID,
	}
	if err := s.repo.Create(ctx, s.db, rule); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, s.db, rule.ID, tenantID)
}

func (s *pricingRuleService) GetByID(ctx context.Context, id, tenantID string) (*entity.PricingRule, error) {
	rule, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "pricing rule not found")
		}
		return nil, err
	}
	return rule, nil
}

func (s *pricingRuleService) List(ctx context.Context, tenantID string) ([]entity.PricingRule, error) {
	return s.repo.List(ctx, s.db, tenantID)
}

func (s *pricingRuleService) Update(ctx context.Context, id, tenantID, actorID string, req request.UpdatePricingRule) (*entity.PricingRule, error) {
	rule := &entity.PricingRule{
		ID:        id,
		TenantID:  tenantID,
		Name:      derefStr(req.Name),
		Value:     derefFloat(req.Value),
		ValidFrom: req.ValidFrom,
		ValidTo:   req.ValidTo,
		UpdatedBy: &actorID,
	}
	if err := s.repo.Update(ctx, s.db, rule); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "pricing rule not found")
		}
		return nil, err
	}
	return s.repo.FindByID(ctx, s.db, id, tenantID)
}

func (s *pricingRuleService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.repo.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "pricing rule not found")
		}
		return err
	}
	return nil
}

// ============================================================
// couponService
// ============================================================

type couponService struct {
	db   *sqlx.DB
	repo repository.CouponRepository
}

func NewCouponService(db *sqlx.DB, repo repository.CouponRepository) CouponService {
	return &couponService{db: db, repo: repo}
}

func (s *couponService) Create(ctx context.Context, tenantID, actorID string, req request.CreateCoupon) (*entity.Coupon, error) {
	// Guard: code must be unique within the tenant.
	if existing, err := s.repo.FindByCode(ctx, s.db, req.Code, tenantID); err == nil && existing != nil {
		return nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "coupon code already exists")
	}

	c := &entity.Coupon{
		ID:            uuid.NewString(),
		TenantID:      tenantID,
		Code:          req.Code,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		MinOrderValue: req.MinOrderValue,
		UsageLimit:    req.UsageLimit,
		ValidFrom:     req.ValidFrom,
		ValidTo:       req.ValidTo,
		CreatedBy:     &actorID,
	}
	if err := s.repo.Create(ctx, s.db, c); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, s.db, c.ID, tenantID)
}

func (s *couponService) GetByID(ctx context.Context, id, tenantID string) (*entity.Coupon, error) {
	c, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "coupon not found")
		}
		return nil, err
	}
	return c, nil
}

func (s *couponService) List(ctx context.Context, tenantID string, page, perPage int) ([]entity.Coupon, error) {
	if perPage <= 0 {
		perPage = 20
	}
	if page <= 0 {
		page = 1
	}
	return s.repo.List(ctx, s.db, tenantID, perPage, (page-1)*perPage)
}

func (s *couponService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.repo.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "coupon not found")
		}
		return err
	}
	return nil
}

func (s *couponService) Validate(ctx context.Context, tenantID string, req request.ValidateCoupon) (*response.CouponValidation, error) {
	c, err := s.repo.FindByCode(ctx, s.db, req.Code, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &response.CouponValidation{Valid: false}, nil
		}
		return nil, err
	}

	now := time.Now()

	// Expiry check.
	if c.ValidTo != nil && c.ValidTo.Before(now) {
		return &response.CouponValidation{Valid: false}, nil
	}
	if c.ValidFrom != nil && c.ValidFrom.After(now) {
		return &response.CouponValidation{Valid: false}, nil
	}

	// Usage limit check.
	if c.UsageLimit != nil && c.UsedCount >= *c.UsageLimit {
		return &response.CouponValidation{Valid: false}, nil
	}

	// Minimum order value check.
	if req.Subtotal < c.MinOrderValue {
		return &response.CouponValidation{Valid: false}, nil
	}

	discountAmount := computeDiscount(c.DiscountType, c.DiscountValue, req.Subtotal)

	return &response.CouponValidation{
		Valid:          true,
		CouponID:       &c.ID,
		DiscountType:   c.DiscountType,
		DiscountValue:  c.DiscountValue,
		DiscountAmount: discountAmount,
	}, nil
}

// ============================================================
// promotionService
// ============================================================

type promotionService struct {
	db   *sqlx.DB
	repo repository.PromotionRepository
}

func NewPromotionService(db *sqlx.DB, repo repository.PromotionRepository) PromotionService {
	return &promotionService{db: db, repo: repo}
}

func (s *promotionService) Create(ctx context.Context, tenantID, actorID string, req request.CreatePromotion) (*entity.Promotion, error) {
	p := &entity.Promotion{
		ID:            uuid.NewString(),
		TenantID:      tenantID,
		Name:          req.Name,
		Description:   req.Description,
		PromotionType: req.PromotionType,
		Value:         req.Value,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		CreatedBy:     &actorID,
	}
	if err := s.repo.Create(ctx, s.db, p); err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, s.db, p.ID, tenantID)
}

func (s *promotionService) GetByID(ctx context.Context, id, tenantID string) (*entity.Promotion, error) {
	p, err := s.repo.FindByID(ctx, s.db, id, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "promotion not found")
		}
		return nil, err
	}
	return p, nil
}

func (s *promotionService) ListActive(ctx context.Context, tenantID string) ([]entity.Promotion, error) {
	return s.repo.ListActive(ctx, s.db, tenantID)
}

func (s *promotionService) Delete(ctx context.Context, id, tenantID string) error {
	if err := s.repo.Delete(ctx, s.db, id, tenantID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return pkgresponse.NewAppError(pkgresponse.CodeNotFound, "promotion not found")
		}
		return err
	}
	return nil
}

// ============================================================
// pricingService — implements PricingQuoter (satisfies booking.PricingQuoter)
// ============================================================

type pricingService struct {
	db      *sqlx.DB
	rules   repository.PricingRuleRepository
	coupons repository.CouponRepository
}

// NewPricingService builds the PricingQuoter consumed by the booking module.
func NewPricingService(db *sqlx.DB, rules repository.PricingRuleRepository, coupons repository.CouponRepository) PricingQuoter {
	return &pricingService{db: db, rules: rules, coupons: coupons}
}

// QuoteItem resolves the best applicable pricing rule for the given
// asset and computes the unit price over the requested duration.
// Falls back to 0 when no rule is found, so bookings are never blocked
// by missing pricing configuration.
func (s *pricingService) QuoteItem(ctx context.Context, tenantID, assetID string, start, end time.Time, qty int) (float64, error) {
	// We need the asset's category to do category-level fallback.
	// The category comes from the inventory module but the pricing
	// module depends only on the repository, not on inventory directly.
	// We pass an empty string for categoryID here; if the caller needs
	// category-level pricing it can populate assetID's associated
	// categoryID via the applicable rule query (which accepts both).
	rule, err := s.rules.FindApplicable(ctx, s.db, tenantID, assetID, "")
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// No rule configured → zero price (operator must configure rules).
			return 0, nil
		}
		return 0, err
	}

	duration := end.Sub(start)
	unitPrice := computeUnitPrice(rule, duration)
	return unitPrice, nil
}

// ValidateCoupon validates a coupon code and returns the discount amount
// plus the coupon ID for storage on the booking. Returning (0, nil, nil)
// when code is empty makes it safe to call unconditionally.
func (s *pricingService) ValidateCoupon(ctx context.Context, tenantID, code string, subtotal float64) (float64, *string, error) {
	if code == "" {
		return 0, nil, nil
	}

	c, err := s.coupons.FindByCode(ctx, s.db, code, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, nil, pkgresponse.NewAppError(pkgresponse.CodeNotFound, "coupon code not found or inactive")
		}
		return 0, nil, err
	}

	now := time.Now()
	if c.ValidTo != nil && c.ValidTo.Before(now) {
		return 0, nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "coupon has expired")
	}
	if c.ValidFrom != nil && c.ValidFrom.After(now) {
		return 0, nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "coupon is not yet valid")
	}
	if c.UsageLimit != nil && c.UsedCount >= *c.UsageLimit {
		return 0, nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "coupon usage limit reached")
	}
	if subtotal < c.MinOrderValue {
		return 0, nil, pkgresponse.NewAppError(pkgresponse.CodeConflict, "order subtotal does not meet the minimum required for this coupon")
	}

	// Increment usage immediately to prevent race condition from
	// concurrent bookings using the same coupon.
	if err := s.coupons.IncrementUsage(ctx, s.db, c.ID); err != nil {
		return 0, nil, err
	}

	discount := computeDiscount(c.DiscountType, c.DiscountValue, subtotal)
	return discount, &c.ID, nil
}

// ============================================================
// pricing helpers
// ============================================================

// computeUnitPrice translates a rule + duration into a price.
func computeUnitPrice(rule *entity.PricingRule, duration time.Duration) float64 {
	switch rule.RuleType {
	case entity.RuleTypeFlat:
		return rule.Value
	case entity.RuleTypePerHour:
		hours := math.Ceil(duration.Hours())
		return rule.Value * hours
	case entity.RuleTypePerDay:
		days := math.Ceil(duration.Hours() / 24)
		return rule.Value * days
	case entity.RuleTypePerWeek:
		weeks := math.Ceil(duration.Hours() / (24 * 7))
		return rule.Value * weeks
	default:
		return rule.Value
	}
}

// computeDiscount applies a coupon's discount type to the subtotal.
func computeDiscount(discountType string, value, subtotal float64) float64 {
	switch discountType {
	case entity.DiscountTypePercentage:
		discount := subtotal * (value / 100)
		if discount > subtotal {
			return subtotal
		}
		return math.Round(discount*100) / 100
	case entity.DiscountTypeFixed:
		if value > subtotal {
			return subtotal
		}
		return value
	default:
		return 0
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
