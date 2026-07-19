// Package validator wraps go-playground/validator behind a small, stable
// surface so call sites never import the underlying library directly.
// Domain-specific custom rules (date ranges, tenant-scoped uniqueness,
// etc.) are registered centrally here in later phases without changing
// this package's exported API.
package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate is the application's request validator.
type Validate struct {
	v *validator.Validate
}

// FieldError describes a single failed validation rule in API-friendly shape.
type FieldError struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
}

// New builds a validator instance with struct-tag validation enabled.
func New() *Validate {
	return &Validate{v: validator.New(validator.WithRequiredStructEnabled())}
}

// Struct validates s and returns a flat list of field errors, or nil if s
// is valid.
func (val *Validate) Struct(s any) []FieldError {
	err := val.v.Struct(s)
	if err == nil {
		return nil
	}

	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return []FieldError{{Field: "_", Rule: "invalid"}}
	}

	out := make([]FieldError, 0, len(verrs))
	for _, fe := range verrs {
		out = append(out, FieldError{
			Field: strings.ToLower(fe.Field()),
			Rule:  fe.Tag(),
		})
	}
	return out
}
