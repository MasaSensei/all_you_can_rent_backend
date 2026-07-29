package service

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateBookingNumber(t *testing.T) {
	n1 := generateBookingNumber()
	n2 := generateBookingNumber()

	if !strings.HasPrefix(n1, "BK-") {
		t.Errorf("booking number should start with BK-, got %s", n1)
	}
	if len(n1) != 11 { // "BK-" (3) + 8 uuid chars = 11
		t.Errorf("booking number should be 11 chars, got %d: %s", len(n1), n1)
	}
	if n1 == n2 {
		t.Error("generateBookingNumber should produce unique values")
	}
}

func TestPassthroughPricer_QuoteItem(t *testing.T) {
	// PassthroughPricer must always return zero so bookings compile
	// end-to-end before the real PricingService is injected.
	p := &PassthroughPricer{}
	now := time.Now()

	price, err := p.QuoteItem(nil, "tenant1", "asset1", now, now.Add(24*time.Hour), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if price != 0 {
		t.Errorf("want 0, got %.2f", price)
	}
}

func TestPassthroughPricer_ValidateCoupon(t *testing.T) {
	p := &PassthroughPricer{}

	discount, couponID, err := p.ValidateCoupon(nil, "tenant1", "CODE10", 100_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if discount != 0 {
		t.Errorf("want discount=0, got %.2f", discount)
	}
	if couponID != nil {
		t.Errorf("want couponID=nil, got %v", couponID)
	}
}
