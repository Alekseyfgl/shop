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

type charDefaultValueHandler struct{}

type CharDefaultValueHandlerInterface interface {
	GetAllDefValue(c *fiber.Ctx) error
	GetDefValueById(c *fiber.Ctx) error
	CreateDefValue(c *fiber.Ctx) error
	UpdateDefValue(c *fiber.Ctx) error
	DeleteDefValue(c *fiber.Ctx) error
}

func NewCharDefaultValueHandler() CharDefaultValueHandlerInterface {
	return &charDefaultValueHandler{}
}

var CharDefaultValueHandler = NewCharDefaultValueHandler()

func (h *charDefaultValueHandler) GetAllDefValue(c *fiber.Ctx) error {
	pageNumberInterface := c.Locals("pageNumber")
	pageSizeInterface := c.Locals("pageSize")

	pageNumber, ok := pageNumberInterface.(int)
	if !ok {
		log.Error("Failed to retrieve page number from context")
		pageNumber = 1
	}

	pageSize, ok := pageSizeInterface.(int)
	if !ok {
		log.Error("Failed to retrieve page char_default_value from context")
		pageSize = 100
	}

	users, err := service.CharDefaultValueService.GetAllDefValue(pageNumber, pageSize)
	if err != nil {
		log.Error("Failed to fetch paginated char_default_value", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch char_default_value", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *charDefaultValueHandler) CreateDefValue(c *fiber.Ctx) error {
	reqInterface := c.Locals("validatedBody")

	body, ok := reqInterface.(dto.CreateCharDefValueRequest)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	size, err := service.CharDefaultValueService.CreateDefValue(&body)
	if err != nil {
		log.Error("Failed to create char_default_value", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create char_default_value", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(size)
}

func (h *charDefaultValueHandler) UpdateDefValue(c *fiber.Ctx) error {
	reqDto := c.Locals("validatedBody")

	body, ok := reqDto.(dto.UpdateCharDefValueRequest)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	size, err := service.CharDefaultValueService.UpdateDefValue(&body)
	if err != nil {
		log.Error("Failed to create char_default_value", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create char_default_value", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(size)
}

func (h *charDefaultValueHandler) DeleteDefValue(c *fiber.Ctx) error {
	idContext := c.Locals("Id")

	IdStr, ok := idContext.(string)
	if !ok || IdStr == "0" {
		log.Error("ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "ID is required", nil).Send(c)
	}

	id, err := utils.StringToInt(IdStr)
	if err != nil {
		return http_error.NewHTTPError(fiber.StatusInternalServerError, err.Error(), nil).Send(c)
	}

	err = service.CharDefaultValueService.DeleteDefValueById(id)
	if err != nil {
		log.Error("Failed to remove char_default_value", zap.Int("id", id), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to remove char_default_value", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": id})
}

func (h *charDefaultValueHandler) GetDefValueById(c *fiber.Ctx) error {
	idStr, ok := c.Locals("Id").(string)
	if !ok || idStr == "0" {
		log.Error("ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "ID is required", nil).Send(c)
	}

	defValId, err := utils.StringToInt(idStr)
	if err != nil {
		return http_error.NewHTTPError(fiber.StatusInternalServerError, err.Error(), nil).Send(c)
	}

	card, err := service.CharDefaultValueService.GetDefValueById(defValId)
	if err != nil {
		log.Error("Failed to find char_default_value", zap.Int("defValId", defValId), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to find card", nil).Send(c)
	}

	return c.JSON(card)
}
