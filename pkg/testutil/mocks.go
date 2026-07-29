// Package testutil provides lightweight helpers and mock implementations
// for unit-testing RentOS modules without a real database or external services.
package testutil

import (
	"context"
	"time"
)

// ---- Mock: InventoryChecker ----

// MockInventoryChecker satisfies booking/service.InventoryChecker.
type MockInventoryChecker struct {
	// Available controls the return value; default true.
	Available bool
	Err       error
}

func (m *MockInventoryChecker) CheckAvailability(_ context.Context, _ string, _, _ time.Time) (bool, error) {
	if m.Err != nil {
		return false, m.Err
	}
	return m.Available, nil
}

// ---- Mock: PricingQuoter ----

// MockPricingQuoter satisfies booking/service.PricingQuoter.
type MockPricingQuoter struct {
	UnitPrice float64
	Discount  float64
	CouponID  *string
	Err       error
}

func (m *MockPricingQuoter) QuoteItem(_ context.Context, _, _ string, _, _ time.Time, _ int) (float64, error) {
	return m.UnitPrice, m.Err
}

func (m *MockPricingQuoter) ValidateCoupon(_ context.Context, _, _ string, _ float64) (float64, *string, error) {
	return m.Discount, m.CouponID, m.Err
}

// ---- Helpers ----

// Ptr returns a pointer to v. Useful for optional fields in test data.
func Ptr[T any](v T) *T { return &v }

// MustTime parses an RFC3339 string and panics on error. For test fixtures only.
func MustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic("testutil.MustTime: " + err.Error())
	}
	return t
}

// DaysFromNow returns a time.Time exactly n days from now (UTC).
func DaysFromNow(n int) time.Time {
	return time.Now().UTC().Add(time.Duration(n) * 24 * time.Hour)
}
