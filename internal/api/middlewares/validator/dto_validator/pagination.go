package dto_validator

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"shop/pkg/http_error"
	"shop/pkg/log"
	"strconv"
)

func ValidatePaginationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// validate 'page' parameter
		pageStr := c.Query("page", "1")
		pageNumber, err := strconv.Atoi(pageStr)
		if err != nil || pageNumber < 1 {
			log.Error("Invalid page number", zap.String("page", pageStr))
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid page number", nil).Send(c)
		}

		// validate 'size' parameter
		sizeStr := c.Query("size", "10")
		pageSize, err := strconv.Atoi(sizeStr)
		if err != nil || pageSize < 1 {
			log.Error("Invalid page size", zap.String("size", sizeStr))
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid page size", nil).Send(c)
		}

		// Store pagination parameters in context
		c.Locals("pageNumber", pageNumber)
		c.Locals("pageSize", pageSize)

		return c.Next()
	}
}
