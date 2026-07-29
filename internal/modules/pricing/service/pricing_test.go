package service

import (
	"math"
	"testing"
	"time"

	"rentos-backend/internal/modules/pricing/entity"
)

// ---- computeUnitPrice ----

func TestComputeUnitPrice_Flat(t *testing.T) {
	rule := &entity.PricingRule{RuleType: entity.RuleTypeFlat, Value: 500_000}
	duration := 3 * 24 * time.Hour
	got := computeUnitPrice(rule, duration)
	if got != 500_000 {
		t.Errorf("flat: want 500000, got %.2f", got)
	}
}

func TestComputeUnitPrice_PerDay(t *testing.T) {
	rule := &entity.PricingRule{RuleType: entity.RuleTypePerDay, Value: 100_000}
	duration := 3 * 24 * time.Hour // exactly 3 days
	want := 300_000.0
	got := computeUnitPrice(rule, duration)
	if got != want {
		t.Errorf("per_day 3d: want %.2f, got %.2f", want, got)
	}
}

func TestComputeUnitPrice_PerDay_Partial(t *testing.T) {
	rule := &entity.PricingRule{RuleType: entity.RuleTypePerDay, Value: 100_000}
	// 1 day + 1 hour → ceil to 2 days
	duration := 25 * time.Hour
	want := 200_000.0
	got := computeUnitPrice(rule, duration)
	if got != want {
		t.Errorf("per_day partial: want %.2f, got %.2f", want, got)
	}
}

func TestComputeUnitPrice_PerHour(t *testing.T) {
	rule := &entity.PricingRule{RuleType: entity.RuleTypePerHour, Value: 10_000}
	duration := 2*time.Hour + 30*time.Minute // 2.5h → ceil 3h
	want := 30_000.0
	got := computeUnitPrice(rule, duration)
	if got != want {
		t.Errorf("per_hour partial: want %.2f, got %.2f", want, got)
	}
}

func TestComputeUnitPrice_PerWeek(t *testing.T) {
	rule := &entity.PricingRule{RuleType: entity.RuleTypePerWeek, Value: 500_000}
	duration := 8 * 24 * time.Hour // 8 days → ceil 2 weeks
	want := 1_000_000.0
	got := computeUnitPrice(rule, duration)
	if got != want {
		t.Errorf("per_week 8d: want %.2f, got %.2f", want, got)
	}
}

// ---- computeDiscount ----

func TestComputeDiscount_Percentage(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		subtotal float64
		want     float64
	}{
		{"10% of 1_000_000", 10, 1_000_000, 100_000},
		{"50% of 200_000", 50, 200_000, 100_000},
		{"100% of 50_000", 100, 50_000, 50_000},
		{"overflow guard: 110% capped at subtotal", 110, 50_000, 50_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeDiscount(entity.DiscountTypePercentage, tt.value, tt.subtotal)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("want %.2f, got %.2f", tt.want, got)
			}
		})
	}
}

func TestComputeDiscount_Fixed(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		subtotal float64
		want     float64
	}{
		{"fixed 50_000 from 200_000", 50_000, 200_000, 50_000},
		{"overflow guard: 300_000 capped at subtotal 200_000", 300_000, 200_000, 200_000},
		{"zero discount", 0, 100_000, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeDiscount(entity.DiscountTypeFixed, tt.value, tt.subtotal)
			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("want %.2f, got %.2f", tt.want, got)
			}
		})
	}
}

func TestComputeDiscount_UnknownType(t *testing.T) {
	got := computeDiscount("unknown", 100, 500_000)
	if got != 0 {
		t.Errorf("unknown type: want 0, got %.2f", got)
	}
}
