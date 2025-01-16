package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"shop/pkg/http_error"
)

func LimitQueryParamsMiddleware(c *fiber.Ctx) error {
	const maxQueryParams = 30
	if len(c.Queries()) > maxQueryParams {
		return http_error.NewHTTPError(fiber.StatusBadRequest, "Too many query parameters", nil).Send(c)
	}
	return c.Next()
}
