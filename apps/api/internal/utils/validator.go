package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError converts validator.ValidationErrors into a friendly string
func FormatValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			switch fe.Tag() {
			case "required":
				return fmt.Sprintf("%s is required", fe.Field())
			case "email":
				return fmt.Sprintf("%s must be a valid email address", fe.Field())
			case "min":
				return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
			case "oneof":
				return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
			}
			return fmt.Sprintf("%s is invalid (%s)", fe.Field(), fe.Tag())
		}
	}
	// Fallback for JSON unmarshal errors or others
	return "Invalid request payload"
}
