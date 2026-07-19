package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrMismatch = errors.New("password: hash mismatch")

const defaultCost = 12

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
