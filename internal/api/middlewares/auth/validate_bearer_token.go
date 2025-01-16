package auth

import (
	"github.com/gofiber/fiber/v2"
	"shop/internal/service"
	"shop/pkg/http_error"
	"strings"
)

func JwtAuthMiddleware(jwtService service.JWTServiceInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем заголовок Authorization
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return http_error.NewHTTPError(fiber.StatusUnauthorized, "Authorization header is missing", nil).Send(c)
		}

		// Проверяем, что заголовок начинается с "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return http_error.NewHTTPError(fiber.StatusUnauthorized, "Invalid authorization format", nil).Send(c)
		}

		// Извлекаем токен из заголовка
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Валидируем токен
		token, err := jwtService.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			return http_error.NewHTTPError(fiber.StatusUnauthorized, "Invalid or expired token", nil).Send(c)
		}

		// Извлекаем claims из токена
		claims, ok := token.Claims.(*service.Claims)
		if !ok || claims.UserId == "" {
			return http_error.NewHTTPError(fiber.StatusUnauthorized, "Invalid token claims", nil).Send(c)
		}

		// Добавляем UserId в локальные данные запроса
		c.Locals("userId", claims.UserId)

		// Продолжаем выполнение запроса
		return c.Next()
	}
}

// GetUserId retrieves the UserId from the Fiber context
func GetUserId(c *fiber.Ctx) (string, bool) {
	userId, ok := c.Locals("userId").(string)
	return userId, ok
}
