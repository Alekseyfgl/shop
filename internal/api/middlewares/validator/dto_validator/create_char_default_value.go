package dto_validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"shop/internal/api/dto"
	"shop/internal/api/middlewares/validator/format_validation_error"
	"shop/pkg/http_error"
	"shop/pkg/log"
)

func ValidateCreateCharDefaultValueMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the input data.
		var req dto.CreateCharDefValueRequest
		if err := c.BodyParser(&req); err != nil {
			log.Error("Failed to parse request body", zap.Error(err))
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid request body", nil).Send(c)
		}

		// validate the input data.
		if err := validate.Struct(&req); err != nil {
			log.Error("Validation failed for request body", zap.Error(err))

			// If the error is a validation error.
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorDetails := format_validation_error.FormatValidationErrors(validationErrors)
				return http_error.NewHTTPError(fiber.StatusBadRequest, "Validation error", errorDetails).Send(c)
			}

			// For other validation errors.
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid input", nil).Send(c)
		}

		// Store the validated data in context for use in the handler.
		c.Locals("validatedBody", req)

		// Proceed to the next handler.
		return c.Next()
	}
}
