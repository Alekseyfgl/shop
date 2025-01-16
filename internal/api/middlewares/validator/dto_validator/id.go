package dto_validator

import (
	"github.com/gofiber/fiber/v2"
	"shop/pkg/http_error"
	"shop/pkg/log"
)

func ValidateIdMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Params("id")
		if userId == "" {
			log.Error("Id param is missing in the request")
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Id param is required", nil).Send(c)
		}

		// Store the ID in context for use in the handler.
		c.Locals("Id", userId)

		return c.Next()
	}
}
