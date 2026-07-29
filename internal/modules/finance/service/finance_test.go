package service

import (
	"math"
	"testing"
)

func TestRound2(t *testing.T) {
	tests := []struct {
		input float64
		want  float64
	}{
		{100.005, 100.01},
		{100.004, 100.00},
		{0.1 + 0.2, 0.30},  // classic float pitfall
		{1_234_567.891, 1_234_567.89},
		{0, 0},
		{-50.555, -50.56},
	}
	for _, tt := range tests {
		got := round2(tt.input)
		if math.Abs(got-tt.want) > 0.001 {
			t.Errorf("round2(%.6f): want %.2f, got %.2f", tt.input, tt.want, got)
		}
	}
}

func TestNormPage(t *testing.T) {
	tests := []struct {
		perPage, page int
		wantPerPage   int
		wantPage      int
	}{
		{20, 1, 20, 1},
		{0, 0, 20, 1},   // defaults
		{-1, -5, 20, 1}, // negatives
		{200, 1, 100, 1}, // cap at 100
		{10, 5, 10, 5},
	}
	for _, tt := range tests {
		pp, p := normPage(tt.perPage, tt.page)
		if pp != tt.wantPerPage || p != tt.wantPage {
			t.Errorf("normPage(%d,%d): want (%d,%d), got (%d,%d)",
				tt.perPage, tt.page, tt.wantPerPage, tt.wantPage, pp, p)
		}
	}
}
