package testutil

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestPtr(t *testing.T) {
	s := "hello"
	p := Ptr(s)
	if p == nil {
		t.Fatal("Ptr returned nil")
	}
	if *p != s {
		t.Errorf("*Ptr(%q) = %q, want %q", s, *p, s)
	}

	n := 42
	pn := Ptr(n)
	if *pn != 42 {
		t.Errorf("*Ptr(42) = %d", *pn)
	}
}

func TestDaysFromNow(t *testing.T) {
	d := DaysFromNow(5)
	diff := d.Sub(time.Now())
	// Should be ~5 days ± a second.
	if diff < 4*24*time.Hour || diff > 6*24*time.Hour {
		t.Errorf("DaysFromNow(5): got %v, want ~5 days", diff)
	}
}

func TestBookingDates(t *testing.T) {
	start, end := BookingDates(1, 3)
	if !end.After(start) {
		t.Error("end should be after start")
	}
	diff := end.Sub(start)
	if diff != 3*24*time.Hour {
		t.Errorf("duration: want 72h, got %v", diff)
	}
}

func TestUniqueEmail(t *testing.T) {
	e1 := UniqueEmail()
	e2 := UniqueEmail()
	if e1 == e2 {
		t.Error("UniqueEmail should produce unique values")
	}
	if !strings.HasSuffix(e1, "@rentos.test") {
		t.Errorf("email should end with @rentos.test, got %s", e1)
	}
}

func TestAssertEqual(t *testing.T) {
	// Test that matching values don't error.
	mock := &mockT{}
	AssertEqual(mock, "expected", "expected", "label")
	if mock.failed {
		t.Error("AssertEqual should not fail for matching values")
	}

	// Test that mismatching values error.
	mock2 := &mockT{}
	AssertEqual(mock2, "want", "got", "label")
	if !mock2.failed {
		t.Error("AssertEqual should fail for mismatching values")
	}
}

func TestAssertErrorContains(t *testing.T) {
	mock := &mockT{}
	AssertErrorContains(mock, errors.New("record not found"), "not found")
	if mock.failed {
		t.Error("should not fail when error contains substr")
	}

	mock2 := &mockT{}
	AssertErrorContains(mock2, nil, "anything")
	if !mock2.failed {
		t.Error("should fail when err is nil")
	}
}

// mockT is a minimal testing.T substitute for testing testutil helpers.
type mockT struct{ failed bool }

func (m *mockT) Helper()                          {}
func (m *mockT) Errorf(_ string, _ ...any)        { m.failed = true }
func (m *mockT) Fatalf(_ string, _ ...any)        { m.failed = true }
