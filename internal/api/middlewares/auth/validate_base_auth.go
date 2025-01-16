package auth

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"shop/configs/env"
	"shop/pkg/http_error"
	"strings"
)

func BasicAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем заголовок Authorization
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return http_error.NewHTTPError(fiber.StatusUnauthorized, "Invalid authorization", nil).Send(c)
		}

		// Проверяем что заголовок начинается с "Basic "
		if !strings.HasPrefix(authHeader, "Basic ") {
			return http_error.NewHTTPError(fiber.StatusUnauthorized, "Invalid authorization", nil).Send(c)
		}

		// Декодируем base64 часть заголовка
		encodedCreds := strings.TrimPrefix(authHeader, "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(encodedCreds)
		if err != nil {
			return http_error.NewHTTPError(fiber.StatusUnauthorized, "Invalid authorization", nil).Send(c)
		}

		// Формат после декодирования: username:password
		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			return http_error.NewHTTPError(fiber.StatusUnauthorized, "Invalid authorization", nil).Send(c)
		}

		username, password := parts[0], parts[1]

		validSuperLogin := env.GetEnv("SUPER_ADMIN_LOGIN", "")
		validSuperPassword := env.GetEnv("SUPER_ADMIN_PASSWORD", "")

		// Проверяем креденшалы
		if username != validSuperLogin || password != validSuperPassword {
			return http_error.NewHTTPError(fiber.StatusUnauthorized, "Invalid authorization", nil).Send(c)
		}

		// Если всё ок — продолжаем
		return c.Next()
	}
}
