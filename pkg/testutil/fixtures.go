package testutil

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ---- Common fixture factories ----
// These produce minimal valid structs suitable for unit tests.
// They do NOT hit the database — pair with mock repositories.

// NewTenantID returns a deterministic tenant UUID for tests.
func NewTenantID() string { return "11111111-1111-1111-1111-111111111111" }

// NewActorID returns a deterministic actor UUID for tests.
func NewActorID() string { return "22222222-2222-2222-2222-222222222222" }

// NewID returns a random UUID string.
func NewID() string { return uuid.NewString() }

// BookingDates returns a start/end pair n days from now.
func BookingDates(startOffsetDays, durationDays int) (start, end time.Time) {
	now := time.Now().UTC().Truncate(24 * time.Hour)
	start = now.Add(time.Duration(startOffsetDays) * 24 * time.Hour)
	end = start.Add(time.Duration(durationDays) * 24 * time.Hour)
	return
}

// ---- Error matchers ----

// AssertErrorContains fails the test if err is nil or does not contain substr.
func AssertErrorContains(t interface {
	Helper()
	Errorf(string, ...any)
}, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error containing %q, got nil", substr)
		return
	}
	if !contains(err.Error(), substr) {
		t.Errorf("error %q does not contain %q", err.Error(), substr)
	}
}

// AssertNoError fails the test if err is not nil.
func AssertNoError(t interface {
	Helper()
	Fatalf(string, ...any)
}, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertEqual fails the test if got != want (formatted with %v).
func AssertEqual[T comparable](t interface {
	Helper()
	Errorf(string, ...any)
}, want, got T, label string) {
	t.Helper()
	if want != got {
		t.Errorf("%s: want %v, got %v", label, want, got)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}

// ---- String helpers ----

// UniqueEmail generates a unique email for test isolation.
func UniqueEmail() string {
	return fmt.Sprintf("test-%s@rentos.test", uuid.NewString()[:8])
}
