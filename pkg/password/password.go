// Package password provides bcrypt hashing and verification.
// Using bcrypt (via golang.org/x/crypto) which is battle-tested and
// available without extra dependencies. The cost factor is configurable
// so it can be increased as hardware improves without changing call sites.
package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// ErrMismatch is returned by Verify when the password does not match
// the stored hash.
var ErrMismatch = errors.New("password: hash mismatch")

const defaultCost = 12

// Hash returns a bcrypt hash of the plaintext password.
func Hash(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), defaultCost)
	if err != nil {
		return "", fmt.Errorf("password: hash: %w", err)
	}
	return string(b), nil
}

// Verify checks plain against the stored bcrypt hash.
// Returns ErrMismatch if they do not match, or another error on failure.
func Verify(plain, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return ErrMismatch
	}
	if err != nil {
		return fmt.Errorf("password: verify: %w", err)
	}
	return nil
}
