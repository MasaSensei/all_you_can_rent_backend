package request

import "time"

// CreatePricingRule creates a new pricing rule.
type CreatePricingRule struct {
	CategoryID   *string    `json:"category_id" validate:"omitempty,uuid"`
	AssetID      *string    `json:"asset_id" validate:"omitempty,uuid"`
	Name         string     `json:"name" validate:"required,max=150"`
	RuleType     string     `json:"rule_type" validate:"required,oneof=flat per_day per_hour per_week"`
	Value        float64    `json:"value" validate:"required,min=0"`
	DurationUnit string     `json:"duration_unit" validate:"required,oneof=hour day week"`
	MinDuration  *int       `json:"min_duration" validate:"omitempty,min=1"`
	MaxDuration  *int       `json:"max_duration" validate:"omitempty,min=1"`
	ValidFrom    *time.Time `json:"valid_from"`
	ValidTo      *time.Time `json:"valid_to"`
}

// UpdatePricingRule updates an existing pricing rule.
type UpdatePricingRule struct {
	Name      *string    `json:"name" validate:"omitempty,max=150"`
	Value     *float64   `json:"value" validate:"omitempty,min=0"`
	ValidFrom *time.Time `json:"valid_from"`
	ValidTo   *time.Time `json:"valid_to"`
}

// CreateCoupon creates a new coupon.
type CreateCoupon struct {
	Code          string     `json:"code" validate:"required,max=50,alphanum"`
	DiscountType  string     `json:"discount_type" validate:"required,oneof=percentage fixed"`
	DiscountValue float64    `json:"discount_value" validate:"required,min=0"`
	MinOrderValue float64    `json:"min_order_value" validate:"min=0"`
	UsageLimit    *int       `json:"usage_limit" validate:"omitempty,min=1"`
	ValidFrom     *time.Time `json:"valid_from"`
	ValidTo       *time.Time `json:"valid_to"`
}

// ValidateCoupon checks a coupon code against an order subtotal.
type ValidateCoupon struct {
	Code     string  `json:"code" validate:"required"`
	Subtotal float64 `json:"subtotal" validate:"required,min=0"`
}

// CreatePromotion creates a new promotion.
type CreatePromotion struct {
	Name          string     `json:"name" validate:"required,max=150"`
	Description   *string    `json:"description" validate:"omitempty,max=500"`
	PromotionType string     `json:"promotion_type" validate:"required,oneof=percentage fixed"`
	Value         float64    `json:"value" validate:"required,min=0"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
}
