package validator

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	nepaliPhoneRegex   = regexp.MustCompile(`^(?:\+977[- ]?)?(?:9[78]\d{8}|98\d{8}|97\d{8}|9[78]\d{8}|01[0-9]{7}|0[1-9][0-9]{6,7})$`)
	nepaliNameRegex    = regexp.MustCompile(`^[\p{L}\s\.\-']+$`)
	allowedCharacters = regexp.MustCompile(`^[a-zA-Z0-9\s\-_@.,'()]+$`)
)

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func phoneValidation(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return nepaliPhoneRegex.MatchString(phone)
}

func ritualDateValidation(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	return date.After(time.Now().Add(-24 * time.Hour))
}

func ritualTimeValidation(fl validator.FieldLevel) bool {
	timeStr := fl.Field().String()
	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		_, err = time.Parse("15:04:05", timeStr)
	}
	return err == nil
}

func amountValidation(fl validator.FieldLevel) bool {
	amount := fl.Field().Float()
	return amount > 0 && amount <= 1000000
}

func nepaliNameValidation(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	return nepaliNameRegex.MatchString(strings.TrimSpace(name))
}

func safeTextValidation(fl validator.FieldLevel) bool {
	text := fl.Field().String()
	return allowedCharacters.MatchString(text)
}

func InitializeValidator() *validator.Validate {
	v := validator.New()

	v.RegisterValidation("nepali_phone", phoneValidation)
	v.RegisterValidation("ritual_date", ritualDateValidation)
	v.RegisterValidation("ritual_time", ritualTimeValidation)
	v.RegisterValidation("amount", amountValidation)
	v.RegisterValidation("nepali_name", nepaliNameValidation)
	v.RegisterValidation("safe_text", safeTextValidation)

	return v
}

func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			ve := ValidationError{
				Field: strings.ToLower(e.Field()),
				Tag:   e.Tag(),
				Value: fmt.Sprintf("%v", e.Value()),
			}
			ve.Message = getValidationMessage(ve)
			errors = append(errors, ve)
		}
	}

	return errors
}

func getValidationMessage(e ValidationError) string {
	messages := map[string]map[string]string{
		"required": {
			"default": "This field is required",
		},
		"email": {
			"default": "Invalid email format",
		},
		"min": {
			"default": "Value is too short",
		},
		"max": {
			"default": "Value is too long",
		},
		"nepali_phone": {
			"default": "Invalid Nepali phone number format",
		},
		"ritual_date": {
			"default": "Ritual date must be in the future",
		},
		"ritual_time": {
			"default": "Invalid time format (use HH:MM or HH:MM:SS)",
		},
		"amount": {
			"default": "Amount must be between 1 and 1,000,000",
		},
		"nepali_name": {
			"default": "Name contains invalid characters",
		},
		"safe_text": {
			"default": "Text contains unsafe characters",
		},
		"gte": {
			"default": "Value must be greater than or equal to the minimum",
		},
		"lte": {
			"default": "Value must be less than or equal to the maximum",
		},
		"len": {
			"default": "Invalid length",
		},
	}

	if tagMessages, ok := messages[e.Tag]; ok {
		if msg, ok := tagMessages[e.Field]; ok {
			return msg
		}
		if msg, ok := tagMessages["default"]; ok {
			return msg
		}
	}

	return fmt.Sprintf("Validation failed on field '%s' with tag '%s'", e.Field, e.Tag)
}
