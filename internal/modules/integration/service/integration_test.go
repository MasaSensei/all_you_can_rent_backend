package service

import (
	"strings"
	"testing"
)

func TestGenerateAPIKey_Format(t *testing.T) {
	raw, prefix, hash, err := generateAPIKey()
	if err != nil {
		t.Fatalf("generateAPIKey error: %v", err)
	}

	if !strings.HasPrefix(raw, "rnt_") {
		t.Errorf("raw key should start with rnt_, got %s", raw[:8])
	}
	if len(raw) != 68 { // "rnt_" (4) + hex(32 bytes) = 4+64 = 68
		t.Errorf("raw key length: want 68, got %d", len(raw))
	}
	if prefix != raw[:12] {
		t.Errorf("prefix should be first 12 chars of raw key")
	}
	if len(hash) != 64 { // SHA-256 hex = 64 chars
		t.Errorf("hash length: want 64, got %d", len(hash))
	}
}

func TestGenerateAPIKey_Uniqueness(t *testing.T) {
	raw1, _, hash1, _ := generateAPIKey()
	raw2, _, hash2, _ := generateAPIKey()

	if raw1 == raw2 {
		t.Error("generateAPIKey should produce unique raw keys")
	}
	if hash1 == hash2 {
		t.Error("generateAPIKey should produce unique hashes")
	}
}

func TestHashKey_Deterministic(t *testing.T) {
	key := "rnt_abc123"
	h1 := hashKey(key)
	h2 := hashKey(key)

	if h1 != h2 {
		t.Error("hashKey should be deterministic for same input")
	}
	if len(h1) != 64 {
		t.Errorf("hashKey length: want 64, got %d", len(h1))
	}
}

func TestHashKey_DifferentInputs(t *testing.T) {
	h1 := hashKey("rnt_key_one")
	h2 := hashKey("rnt_key_two")

	if h1 == h2 {
		t.Error("different keys should produce different hashes")
	}
}

func TestResolveByRawKey_HashMatch(t *testing.T) {
	// Verify the hash produced by generateAPIKey matches what hashKey produces
	// for the same raw value — this is the critical invariant for key resolution.
	raw, _, storedHash, err := generateAPIKey()
	if err != nil {
		t.Fatalf("generateAPIKey: %v", err)
	}

	resolvedHash := hashKey(raw)
	if resolvedHash != storedHash {
		t.Errorf("hash mismatch: stored=%s, resolved=%s", storedHash, resolvedHash)
	}
}
