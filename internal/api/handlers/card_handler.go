package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"shop/internal/api/dto"
	"shop/internal/model"
	"shop/internal/service"
	"shop/pkg/http_error"
	"shop/pkg/log"
	"shop/pkg/utils"
	"strings"
)

type cardHandler struct{}

type CardHandlerInterface interface {
	GetCardById(c *fiber.Ctx) error
	GetAllCards(c *fiber.Ctx) error
	GetCardsByVector(c *fiber.Ctx) error
	CreateCard(c *fiber.Ctx) error
}

func NewCardHandler() CardHandlerInterface {
	return &cardHandler{}
}

var CardHandler = NewCardHandler()

func (h *cardHandler) GetCardById(c *fiber.Ctx) error {

	// Retrieve article ID from context.
	id := c.Locals("Id")
	cardIdStr, ok := id.(string)
	if !ok || cardIdStr == "0" {
		log.Error("ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "ID is required", nil).Send(c)
	}

	cardId, err := utils.StringToInt(cardIdStr)
	if err != nil {
		return http_error.NewHTTPError(fiber.StatusInternalServerError, err.Error(), nil).Send(c)
	}

	card, err := service.CardService.GetCardById(cardId)

	if err != nil {
		log.Error("Failed to find card", zap.Int("cardId", cardId), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to find card", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(card)
}

func (h *cardHandler) GetAllCards(c *fiber.Ctx) error {
	// Извлечение query-параметров пагинации
	pageNumber, err := utils.StringToInt(c.Query("pageNumber", "1"))
	if err != nil || pageNumber < 1 {
		log.Error("Invalid or missing page number in query", zap.Error(err))
		pageNumber = 1
	}

	pageSize, err := utils.StringToInt(c.Query("pageSize", "50"))
	if err != nil || pageSize < 1 {
		log.Error("Invalid or missing page size in query", zap.Error(err))
		pageSize = 50
	}

	filters := make([]model.CardFilter, 0)

	// c.Queries() вернёт map[string]string, где key — это имя параметра, а value — его значение
	for key, value := range c.Queries() {
		// Пропускаем параметры пагинации
		if key == "pageNumber" || key == "pageSize" {
			continue
		}

		// Проверяем, содержит ли значение запятые
		if strings.Contains(value, ",") {
			parts := strings.Split(value, ",")
			for _, part := range parts {
				part := strings.TrimSpace(part) // убираем пробелы вокруг
				if part != "" {
					// Каждый элемент добавляем отдельным CardFilter
					filters = append(filters, model.CardFilter{
						Key:    key,
						Values: part,
					})
				}
			}
		} else {
			// Если запятых нет — добавляем значение как есть
			filters = append(filters, model.CardFilter{
				Key:    key,
				Values: value,
			})
		}
	}

	// Логирование параметров для отладки
	log.Info("Fetching cards with filters",
		zap.Any("filters", filters),
	)

	// Вызов сервиса для получения карт с учетом фильтров
	cards, err := service.CardService.GetAllCards(pageNumber, pageSize, &filters)
	if err != nil {
		log.Error("Failed to fetch paginated cards with filters", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch cards", nil).Send(c)
	}

	// Возврат результата
	return c.Status(fiber.StatusOK).JSON(cards)
}

func (h *cardHandler) CreateCard(c *fiber.Ctx) error {
	reqInterface := c.Locals("validatedBody")

	body, ok := reqInterface.(dto.CreateCardDTO)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	// 3. Вызываем метод сервиса
	newID, err := service.CardService.CreateCard(&body)
	if err != nil {
		log.Error("Failed to create card", zap.Error(err))
		// При желании можно вернуть подробную ошибку
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 4. Возвращаем успешный результат
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": newID,
	})
}

func (h *cardHandler) GetCardsByVector(c *fiber.Ctx) error {
	reqInterface := c.Locals("validatedBody")

	body, ok := reqInterface.(dto.GetCardsByVectorDTO)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	cards, err := service.CardService.GetCardsByVector(&body)
	if err != nil {
		log.Error("Failed to get cards by vector", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to get cards", nil).Send(c)
	}

	return c.Status(fiber.StatusCreated).JSON(cards)
}
