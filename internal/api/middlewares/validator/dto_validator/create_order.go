package dto_validator

import (
	"shop/internal/api/dto"
	"shop/internal/api/middlewares/validator/format_validation_error"
	"shop/pkg/http_error"
	"shop/pkg/log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ValidateCreateOrderMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req []dto.OrderDTO
		if err := c.BodyParser(&req); err != nil {
			log.Error("Failed to parse request body", zap.Error(err))
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid request body", nil).Send(c)
		}

		// Проверяем что массив не пустой
		if len(req) == 0 {
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Order array cannot be empty", nil).Send(c)
		}

		// validate each item in the array
		for i, order := range req {
			if err := validate.Struct(&order); err != nil {
				log.Error("Validation failed for request body", zap.Error(err), zap.Int("orderIndex", i))

				// If the error is a validation error.
				if validationErrors, ok := err.(validator.ValidationErrors); ok {
					errorDetails := format_validation_error.FormatValidationErrors(validationErrors)
					return http_error.NewHTTPError(fiber.StatusBadRequest, "Validation error", errorDetails).Send(c)
				}

				// For other validation errors.
				return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid input", nil).Send(c)
			}
		}

		// Store the validated data in context for use in the handler.
		c.Locals("validatedBody", req)

		// Proceed to the next handler.
		return c.Next()
	}
}
