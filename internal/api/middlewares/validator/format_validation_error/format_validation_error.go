package format_validation_error

import (
	"github.com/go-playground/validator/v10"
	"shop/pkg/http_error"
)

// FormatValidationErrors formats validation errors into an array of ErrorItem.
func FormatValidationErrors(validationErrors validator.ValidationErrors) []http_error.ErrorItem {
	var errorDetails []http_error.ErrorItem

	// Опционально можно заранее определить сообщения по тегам:
	errorMessages := map[string]string{
		"required":           "This field is required",
		"min":                "The field does not meet the minimum length requirement",
		"img_base64_or_null": "The field must be null or a valid Base64 string",
		"custom_email":       "Invalid email",
	}

	for _, validationErr := range validationErrors {
		field := validationErr.Field()
		tag := validationErr.Tag()

		errorMsg, ok := errorMessages[tag]
		if !ok {
			errorMsg = "Validation failed on tag: " + tag
		}

		errorDetails = append(errorDetails, http_error.ErrorItem{
			Field: field,
			Error: errorMsg,
		})
	}

	return errorDetails
}
