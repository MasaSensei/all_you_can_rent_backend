package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestSignPayload_Deterministic(t *testing.T) {
	secret := "super-secret-key-minimum-16chars"
	payload := `{"event":"booking.confirmed","booking_id":"abc123"}`

	sig1 := signPayload(secret, payload)
	sig2 := signPayload(secret, payload)

	if sig1 != sig2 {
		t.Error("signPayload should be deterministic for same input")
	}
	if len(sig1) != 64 { // HMAC-SHA256 hex = 64 chars
		t.Errorf("signature length: want 64, got %d", len(sig1))
	}
}

func TestSignPayload_DifferentSecrets(t *testing.T) {
	payload := `{"event":"payment.succeeded"}`
	sig1 := signPayload("secret-one-aaaaaaaa", payload)
	sig2 := signPayload("secret-two-bbbbbbbb", payload)

	if sig1 == sig2 {
		t.Error("different secrets should produce different signatures")
	}
}

func TestSignPayload_DifferentPayloads(t *testing.T) {
	secret := "shared-secret-rentos-webhook"
	sig1 := signPayload(secret, `{"event":"booking.created"}`)
	sig2 := signPayload(secret, `{"event":"booking.cancelled"}`)

	if sig1 == sig2 {
		t.Error("different payloads should produce different signatures")
	}
}

func TestSignPayload_CorrectHMAC(t *testing.T) {
	// Verify against a manual HMAC-SHA256 computation.
	secret := "test-secret-rentos"
	payload := "hello rentos"

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))

	got := signPayload(secret, payload)
	if got != expected {
		t.Errorf("HMAC mismatch:\n  want: %s\n   got: %s", expected, got)
	}
}

func TestSignPayload_EmptyPayload(t *testing.T) {
	// Empty payload should still produce a valid signature (not panic).
	sig := signPayload("any-secret-abc123", "")
	if len(sig) != 64 {
		t.Errorf("empty payload signature length: want 64, got %d", len(sig))
	}
}
