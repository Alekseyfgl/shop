package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"shop/internal/api/dto"
	"shop/internal/service"
	"shop/pkg/http_error"
	"shop/pkg/log"
	"shop/pkg/utils"
)

type characteristicHandler struct{}

type CharacteristicHandlerInterface interface {
	GetAllCharacteristics(c *fiber.Ctx) error
	GetCharForFilters(c *fiber.Ctx) error
	CreateCharacteristic(c *fiber.Ctx) error
	DeleteCharacteristic(c *fiber.Ctx) error
	UpdateCharacteristic(c *fiber.Ctx) error
}

func NewCharacteristicHandler() CharacteristicHandlerInterface {
	return &characteristicHandler{}
}

var CharacteristicHandler = NewCharacteristicHandler()

func (h *characteristicHandler) GetAllCharacteristics(c *fiber.Ctx) error {
	pageNumberInterface := c.Locals("pageNumber")
	pageSizeInterface := c.Locals("pageSize")

	pageNumber, ok := pageNumberInterface.(int)
	if !ok {
		log.Error("Failed to retrieve page number from context")
		pageNumber = 1
	}

	pageSize, ok := pageSizeInterface.(int)
	if !ok {
		log.Error("Failed to retrieve page characteristic from context")
		pageSize = 100
	}

	users, err := service.CharacteristicService.GetAllCharacteristic(pageNumber, pageSize)
	if err != nil {
		log.Error("Failed to fetch paginated characteristics", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch characteristic", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *characteristicHandler) GetCharForFilters(c *fiber.Ctx) error {
	filters, err := service.CharacteristicService.GetCharForFilters()
	if err != nil {
		log.Error("Failed to fetch paginated filters", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch filters", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(filters)
}

func (h *characteristicHandler) CreateCharacteristic(c *fiber.Ctx) error {
	reqInterface := c.Locals("validatedBody")

	body, ok := reqInterface.(dto.CreateCharacteristicRequest)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	size, err := service.CharacteristicService.CreateCharacteristic(&body)
	if err != nil {
		log.Error("Failed to create characteristic", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create characteristic", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(size)
}

func (h *characteristicHandler) UpdateCharacteristic(c *fiber.Ctx) error {
	reqDto := c.Locals("validatedBody")

	body, ok := reqDto.(dto.UpdateCharacteristicRequest)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	size, err := service.CharacteristicService.UpdateCharacteristic(&body)
	if err != nil {
		log.Error("Failed to create characteristic", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create characteristic", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(size)
}

func (h *characteristicHandler) DeleteCharacteristic(c *fiber.Ctx) error {

	// Retrieve article ID from context.
	id := c.Locals("Id")
	sizeIdStr, ok := id.(string)
	if !ok || sizeIdStr == "0" {
		log.Error("ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "ID is required", nil).Send(c)
	}

	sizeId, err := utils.StringToInt(sizeIdStr)
	if err != nil {
		return http_error.NewHTTPError(fiber.StatusInternalServerError, err.Error(), nil).Send(c)
	}

	err = service.CharacteristicService.DeleteCharacteristic(sizeId)
	if err != nil {
		log.Error("Failed to remove size", zap.Int("sizeId", sizeId), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to remove user", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": sizeId})
}
