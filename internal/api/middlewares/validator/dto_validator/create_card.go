package dto_validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"shop/internal/api/dto"
	"shop/internal/api/middlewares/validator/format_validation_error"
	"shop/internal/repository"
	"shop/pkg/http_error"
	"shop/pkg/log"
)

func ValidateCreateCardMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the input data.
		var req dto.CreateCardDTO
		if err := c.BodyParser(&req); err != nil {
			log.Error("Failed to parse request body", zap.Error(err))
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid request body", nil).Send(c)
		}

		var ids []int = make([]int, 0, len(req.Characteristics))
		for _, char := range req.Characteristics {
			ids = append(ids, char.Id)
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

		err := repository.CharacteristicRepo.CheckCharsByIds(ids)
		if err != nil {
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid input", []http_error.ErrorItem{{
				Field: "id",
				Error: err.Error(),
			}}).Send(c)
		}
		// Store the validated data in context for use in the handler.
		c.Locals("validatedBody", req)

		// Proceed to the next handler.
		return c.Next()
	}
}
