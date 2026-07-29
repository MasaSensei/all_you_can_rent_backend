package entity

import "time"

// PricingRule mirrors the pricing_rules table.
// Rule resolution priority: asset-level overrides category-level.
type PricingRule struct {
	ID           string     `db:"id"`
	TenantID     string     `db:"tenant_id"`
	CategoryID   *string    `db:"category_id"`
	AssetID      *string    `db:"asset_id"`
	Name         string     `db:"name"`
	RuleType     string     `db:"rule_type"`   // flat, per_day, per_hour, per_week
	Value        float64    `db:"value"`
	DurationUnit string     `db:"duration_unit"` // hour, day, week
	MinDuration  *int       `db:"min_duration"`
	MaxDuration  *int       `db:"max_duration"`
	ValidFrom    *time.Time `db:"valid_from"`
	ValidTo      *time.Time `db:"valid_to"`
	Status       string     `db:"status"`
	CreatedBy    *string    `db:"created_by"`
	UpdatedBy    *string    `db:"updated_by"`
	DeletedBy    *string    `db:"deleted_by"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
	Version      int        `db:"version"`
}

const (
	RuleTypeFlat    = "flat"
	RuleTypePerDay  = "per_day"
	RuleTypePerHour = "per_hour"
	RuleTypePerWeek = "per_week"

	DurationUnitHour = "hour"
	DurationUnitDay  = "day"
	DurationUnitWeek = "week"
)
