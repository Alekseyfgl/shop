package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ErrorHandlerMiddleware handles global errors and logs them
func ErrorHandlerMiddleware(logger *zap.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				// Log the error with the provided logger
				logger.Error("Recovered from panic",
					zap.Any("error", r),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
				)

				// Respond with a 500 status code
				_ = c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}
		}()
		return c.Next()
	}
}
