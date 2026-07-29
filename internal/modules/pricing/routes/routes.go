package routes

import (
	"github.com/gofiber/fiber/v2"

	"rentos-backend/internal/modules/pricing/handler"
)

// Register mounts all pricing routes under /api/v1.
func Register(router fiber.Router, h *handler.Handler) {
	rules := router.Group("/pricing-rules")
	rules.Post("/", h.CreatePricingRule)
	rules.Get("/", h.ListPricingRules)
	rules.Get("/:id", h.GetPricingRule)
	rules.Put("/:id", h.UpdatePricingRule)
	rules.Delete("/:id", h.DeletePricingRule)

	coupons := router.Group("/coupons")
	coupons.Post("/", h.CreateCoupon)
	coupons.Get("/", h.ListCoupons)
	coupons.Get("/:id", h.GetCoupon)
	coupons.Delete("/:id", h.DeleteCoupon)
	coupons.Post("/validate", h.ValidateCoupon)

	promotions := router.Group("/promotions")
	promotions.Post("/", h.CreatePromotion)
	promotions.Get("/", h.ListActivePromotions)
	promotions.Get("/:id", h.GetPromotion)
	promotions.Delete("/:id", h.DeletePromotion)
}
