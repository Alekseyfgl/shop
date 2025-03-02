package dto_validator

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"shop/internal/repository"
	"shop/pkg/http_error"
	"shop/pkg/log"
	"strconv"
)

// ValidateNodeTypeIdMiddleware проверяет корректность и существование nodeTypeId
// в query-параметре и при необходимости прерывает выполнение с соответствующей ошибкой.
func ValidateNodeTypeIdMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {

		if !c.Context().QueryArgs().Has("nodeTypeId") {
			c.Locals("nodeTypeId", 0)
			return c.Next()
		}

		nodeTypeParam := c.Query("nodeTypeId", "")
		if nodeTypeParam == "" {
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid nodeTypeId parameter", nil).Send(c)
		}

		// Преобразуем строку в число
		nodeTypeID, err := strconv.Atoi(nodeTypeParam)
		if err != nil {
			// Если преобразование не удалось — это ошибка клиента
			log.Error("Invalid nodeTypeId parameter", zap.String("nodeTypeId", nodeTypeParam), zap.Error(err))
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid nodeTypeId parameter", nil).Send(c)
		}

		nodeType, err := repository.NodeTypeRepo.GetNodeTypeById(nodeTypeID)
		if err != nil {
			// Если не найдено — возвращаем 404
			log.Error("nodeTypeId not found in repository",
				zap.Int("nodeTypeId", nodeTypeID),
				zap.Error(err))
			return http_error.NewHTTPError(fiber.StatusNotFound, "NodeType not found", nil).Send(c)
		}

		// При необходимости сохраним данные в контекст Fiber (для дальнейших хендлеров)
		c.Locals("nodeTypeId", nodeType.ID)

		return c.Next()
	}
}
