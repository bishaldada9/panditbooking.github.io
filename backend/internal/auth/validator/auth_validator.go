package validator

import (
	"github.com/go-playground/validator/v10"

	"github.com/bees/hindu-ritual-platform/pkg/security"
)

func RegisterAuthValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("password_strength", validatePasswordStrength); err != nil {
		return err
	}
	return nil
}

func validatePasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return security.ValidatePasswordStrength(password) == nil
}
