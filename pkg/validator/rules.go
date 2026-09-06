package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func registerRules(v *validator.Validate) {
	_ = v.RegisterValidation("alphanum_dash", validateAlphaNumDash)
}

func validateAlphaNumDash(fl validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(
		`^[a-z0-9]+(?:-[a-z0-9]+)*$`,
		fl.Field().String(),
	)

	return matched
}
