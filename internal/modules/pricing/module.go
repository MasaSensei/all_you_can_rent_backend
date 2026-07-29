package pricing

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/bootstrap"
	"rentos-backend/internal/modules/pricing/handler"
	"rentos-backend/internal/modules/pricing/repository/postgres"
	"rentos-backend/internal/modules/pricing/routes"
	"rentos-backend/internal/modules/pricing/service"
)

// Module holds the pricing module's wired handler and services.
type Module struct {
	handler    *handler.Handler
	quoter     service.PricingQuoter
}

// New builds the pricing module: repositories → services → handler.
func New(c *bootstrap.Container) *Module {
	ruleRepo := postgres.NewPricingRuleRepository(
		query("create_pricing_rule.sql"),
		query("find_pricing_rule_by_id.sql"),
		query("list_pricing_rules.sql"),
		query("find_applicable_pricing_rule.sql"),
		query("update_pricing_rule.sql"),
		query("delete_pricing_rule.sql"),
	)
	couponRepo := postgres.NewCouponRepository(
		query("create_coupon.sql"),
		query("find_coupon_by_id.sql"),
		query("find_coupon_by_code.sql"),
		query("list_coupons.sql"),
		query("increment_coupon_usage.sql"),
		query("delete_coupon.sql"),
	)
	promotionRepo := postgres.NewPromotionRepository(
		query("create_promotion.sql"),
		query("find_promotion_by_id.sql"),
		query("list_active_promotions.sql"),
		query("delete_promotion.sql"),
	)

	ruleSvc := service.NewPricingRuleService(c.DB, ruleRepo)
	couponSvc := service.NewCouponService(c.DB, couponRepo)
	promotionSvc := service.NewPromotionService(c.DB, promotionRepo)

	// PricingQuoter is the cross-module contract consumed by booking.
	quoter := service.NewPricingService(c.DB, ruleRepo, couponRepo)

	h := handler.New(ruleSvc, couponSvc, promotionSvc, c.Validator)
	return &Module{handler: h, quoter: quoter}
}

// RegisterRoutes mounts the module's routes onto /api/v1.
func (m *Module) RegisterRoutes(router fiber.Router) {
	routes.Register(router, m.handler)
}

// PricingQuoter exposes the quoter so the booking module can receive it
// via its PricingQuoter interface, replacing PassthroughPricer.
func (m *Module) PricingQuoter() service.PricingQuoter {
	return m.quoter
}
