package service

import "testing"

func TestRenderTemplate_BasicSubstitution(t *testing.T) {
	body := "Halo {{name}}, booking {{booking_number}} telah dikonfirmasi."
	data := map[string]string{
		"name":           "Budi",
		"booking_number": "BK-a1b2c3d4",
	}
	want := "Halo Budi, booking BK-a1b2c3d4 telah dikonfirmasi."
	got := renderTemplate(body, data)
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestRenderTemplate_MissingKey(t *testing.T) {
	// Missing keys should be left as-is (not replaced).
	body := "Halo {{name}}, kode OTP: {{otp}}"
	data := map[string]string{"name": "Siti"}
	got := renderTemplate(body, data)
	want := "Halo Siti, kode OTP: {{otp}}"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestRenderTemplate_EmptyData(t *testing.T) {
	body := "Notifikasi sistem."
	got := renderTemplate(body, nil)
	if got != body {
		t.Errorf("empty data should return body unchanged, got %q", got)
	}
}

func TestRenderTemplate_MultipleOccurrences(t *testing.T) {
	body := "{{greeting}} {{name}}! Terima kasih {{name}}."
	data := map[string]string{"greeting": "Halo", "name": "Andi"}
	want := "Halo Andi! Terima kasih Andi."
	got := renderTemplate(body, data)
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestRenderTemplate_EmptyBody(t *testing.T) {
	got := renderTemplate("", map[string]string{"key": "val"})
	if got != "" {
		t.Errorf("empty body should stay empty, got %q", got)
	}
}

func TestDerefStr(t *testing.T) {
	s := "rentos"
	if derefStr(&s) != "rentos" {
		t.Error("derefStr: want 'rentos'")
	}
	if derefStr(nil) != "" {
		t.Error("derefStr(nil): want empty string")
	}
}

func TestNormPage(t *testing.T) {
	tests := []struct {
		perPage, page     int
		wantPP, wantPage  int
	}{
		{20, 1, 20, 1},
		{0, 0, 20, 1},
		{-5, -1, 20, 1},
		{200, 3, 100, 3},
		{50, 7, 50, 7},
	}
	for _, tt := range tests {
		pp, p := normPage(tt.perPage, tt.page)
		if pp != tt.wantPP || p != tt.wantPage {
			t.Errorf("normPage(%d,%d) = (%d,%d), want (%d,%d)",
				tt.perPage, tt.page, pp, p, tt.wantPP, tt.wantPage)
		}
	}
}
