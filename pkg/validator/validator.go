package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validate struct {
	v *validator.Validate
}

type FieldError struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
}

func New() *Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	registerRules(v)

	return &Validate{v: v}
}

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
