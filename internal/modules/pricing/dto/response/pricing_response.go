package response

import "time"

// PricingRule is the API-facing shape of entity.PricingRule.
type PricingRule struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	CategoryID   *string    `json:"category_id,omitempty"`
	AssetID      *string    `json:"asset_id,omitempty"`
	Name         string     `json:"name"`
	RuleType     string     `json:"rule_type"`
	Value        float64    `json:"value"`
	DurationUnit string     `json:"duration_unit"`
	MinDuration  *int       `json:"min_duration,omitempty"`
	MaxDuration  *int       `json:"max_duration,omitempty"`
	ValidFrom    *time.Time `json:"valid_from,omitempty"`
	ValidTo      *time.Time `json:"valid_to,omitempty"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Coupon is the API-facing shape of entity.Coupon.
type Coupon struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	Code          string     `json:"code"`
	DiscountType  string     `json:"discount_type"`
	DiscountValue float64    `json:"discount_value"`
	MinOrderValue float64    `json:"min_order_value"`
	UsageLimit    *int       `json:"usage_limit,omitempty"`
	UsedCount     int        `json:"used_count"`
	ValidFrom     *time.Time `json:"valid_from,omitempty"`
	ValidTo       *time.Time `json:"valid_to,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// CouponValidation is returned by the validate-coupon endpoint.
type CouponValidation struct {
	Valid         bool    `json:"valid"`
	CouponID      *string `json:"coupon_id,omitempty"`
	DiscountType  string  `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
	DiscountAmount float64 `json:"discount_amount"`
}

// Promotion is the API-facing shape of entity.Promotion.
type Promotion struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	PromotionType string     `json:"promotion_type"`
	Value         float64    `json:"value"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
