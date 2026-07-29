package handler

import (
	"github.com/gofiber/fiber/v2"

	pricingreq "rentos-backend/internal/modules/pricing/dto/request"
	"rentos-backend/internal/modules/pricing/entity"
	"rentos-backend/internal/modules/pricing/dto/response"
	"rentos-backend/internal/modules/pricing/service"
	apiresponse "rentos-backend/pkg/response"
	"rentos-backend/pkg/validator"
)

const (
	ctxKeyTenantID = "tenant_id"
	ctxKeyUserID   = "user_id"
)

// Handler groups the pricing module's HTTP handlers.
type Handler struct {
	rules      service.PricingRuleService
	coupons    service.CouponService
	promotions service.PromotionService
	validate   *validator.Validate
}

func New(
	rules service.PricingRuleService,
	coupons service.CouponService,
	promotions service.PromotionService,
	v *validator.Validate,
) *Handler {
	return &Handler{rules: rules, coupons: coupons, promotions: promotions, validate: v}
}

func tenantID(c *fiber.Ctx) string {
	if id, ok := c.Locals(ctxKeyTenantID).(string); ok {
		return id
	}
	return c.Get("X-Tenant-ID")
}

func userID(c *fiber.Ctx) string {
	id, _ := c.Locals(ctxKeyUserID).(string)
	return id
}

// ---- Pricing Rules ----

func (h *Handler) CreatePricingRule(c *fiber.Ctx) error {
	var req pricingreq.CreatePricingRule
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	rule, err := h.rules.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, toPricingRuleResponse(rule))
}

func (h *Handler) GetPricingRule(c *fiber.Ctx) error {
	rule, err := h.rules.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toPricingRuleResponse(rule))
}

func (h *Handler) ListPricingRules(c *fiber.Ctx) error {
	rules, err := h.rules.List(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	out := make([]response.PricingRule, 0, len(rules))
	for _, r := range rules {
		out = append(out, toPricingRuleResponse(&r))
	}
	return apiresponse.Success(c, out)
}

func (h *Handler) UpdatePricingRule(c *fiber.Ctx) error {
	var req pricingreq.UpdatePricingRule
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	rule, err := h.rules.Update(c.Context(), c.Params("id"), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toPricingRuleResponse(rule))
}

func (h *Handler) DeletePricingRule(c *fiber.Ctx) error {
	if err := h.rules.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- Coupons ----

func (h *Handler) CreateCoupon(c *fiber.Ctx) error {
	var req pricingreq.CreateCoupon
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	coupon, err := h.coupons.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, toCouponResponse(coupon))
}

func (h *Handler) GetCoupon(c *fiber.Ctx) error {
	coupon, err := h.coupons.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toCouponResponse(coupon))
}

func (h *Handler) ListCoupons(c *fiber.Ctx) error {
	coupons, err := h.coupons.List(c.Context(), tenantID(c),
		c.QueryInt("page", 1), c.QueryInt("per_page", 20),
	)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	out := make([]response.Coupon, 0, len(coupons))
	for _, coupon := range coupons {
		out = append(out, toCouponResponse(&coupon))
	}
	return apiresponse.Success(c, out)
}

func (h *Handler) DeleteCoupon(c *fiber.Ctx) error {
	if err := h.coupons.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

func (h *Handler) ValidateCoupon(c *fiber.Ctx) error {
	var req pricingreq.ValidateCoupon
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	result, err := h.coupons.Validate(c.Context(), tenantID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, result)
}

// ---- Promotions ----

func (h *Handler) CreatePromotion(c *fiber.Ctx) error {
	var req pricingreq.CreatePromotion
	if err := c.BodyParser(&req); err != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "invalid request body"))
	}
	if errs := h.validate.Struct(req); errs != nil {
		return apiresponse.Error(c, apiresponse.NewAppError(apiresponse.CodeValidation, "validation failed").WithDetails(errs))
	}
	promo, err := h.promotions.Create(c.Context(), tenantID(c), userID(c), req)
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Created(c, toPromotionResponse(promo))
}

func (h *Handler) GetPromotion(c *fiber.Ctx) error {
	promo, err := h.promotions.GetByID(c.Context(), c.Params("id"), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.Success(c, toPromotionResponse(promo))
}

func (h *Handler) ListActivePromotions(c *fiber.Ctx) error {
	promos, err := h.promotions.ListActive(c.Context(), tenantID(c))
	if err != nil {
		return apiresponse.FromError(c, err)
	}
	out := make([]response.Promotion, 0, len(promos))
	for _, p := range promos {
		out = append(out, toPromotionResponse(&p))
	}
	return apiresponse.Success(c, out)
}

func (h *Handler) DeletePromotion(c *fiber.Ctx) error {
	if err := h.promotions.Delete(c.Context(), c.Params("id"), tenantID(c)); err != nil {
		return apiresponse.FromError(c, err)
	}
	return apiresponse.NoContent(c)
}

// ---- mapping ----

func toPricingRuleResponse(r *entity.PricingRule) response.PricingRule {
	return response.PricingRule{
		ID: r.ID, TenantID: r.TenantID, CategoryID: r.CategoryID,
		AssetID: r.AssetID, Name: r.Name, RuleType: r.RuleType, Value: r.Value,
		DurationUnit: r.DurationUnit, MinDuration: r.MinDuration, MaxDuration: r.MaxDuration,
		ValidFrom: r.ValidFrom, ValidTo: r.ValidTo,
		Status: r.Status, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func toCouponResponse(c *entity.Coupon) response.Coupon {
	return response.Coupon{
		ID: c.ID, TenantID: c.TenantID, Code: c.Code,
		DiscountType: c.DiscountType, DiscountValue: c.DiscountValue,
		MinOrderValue: c.MinOrderValue, UsageLimit: c.UsageLimit, UsedCount: c.UsedCount,
		ValidFrom: c.ValidFrom, ValidTo: c.ValidTo,
		Status: c.Status, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt,
	}
}

func toPromotionResponse(p *entity.Promotion) response.Promotion {
	return response.Promotion{
		ID: p.ID, TenantID: p.TenantID, Name: p.Name,
		Description: p.Description, PromotionType: p.PromotionType, Value: p.Value,
		StartDate: p.StartDate, EndDate: p.EndDate,
		Status: p.Status, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}
