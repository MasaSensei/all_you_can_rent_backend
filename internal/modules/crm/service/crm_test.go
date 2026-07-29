package service

import "testing"

func TestNormPage(t *testing.T) {
	tests := []struct {
		perPage, page int
		wantPP, wantP int
	}{
		{20, 1, 20, 1},
		{0, 0, 20, 1},
		{-10, -3, 20, 1},
		{150, 2, 100, 2},
		{5, 10, 5, 10},
	}
	for _, tt := range tests {
		pp, p := normPage(tt.perPage, tt.page)
		if pp != tt.wantPP || p != tt.wantP {
			t.Errorf("normPage(%d,%d) = (%d,%d), want (%d,%d)",
				tt.perPage, tt.page, pp, p, tt.wantPP, tt.wantP)
		}
	}
}

func TestDerefStr(t *testing.T) {
	s := "hello"
	if derefStr(&s) != "hello" {
		t.Error("derefStr: want 'hello'")
	}
	if derefStr(nil) != "" {
		t.Error("derefStr(nil): want empty string")
	}
}
