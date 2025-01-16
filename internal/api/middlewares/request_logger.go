package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RequestLoggerMiddleware logs request processing time along with request and response IDs
func RequestLoggerMiddleware(logger *zap.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Generate a unique request ID
		requestID := uuid.New().String()
		c.Set("X-Request-ID", requestID)

		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		// Get the response status code
		statusCode := c.Response().StatusCode()

		// Log the request and response details
		logger.Info("Request processed",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status_code", statusCode),
			zap.Duration("duration", duration))

		return err
	}
}
