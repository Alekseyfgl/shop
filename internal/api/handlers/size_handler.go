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

type sizeHandler struct{}

type SizeHandlerInterface interface {
	GetAllSizes(c *fiber.Ctx) error
	CreateSize(c *fiber.Ctx) error
	DeleteSize(c *fiber.Ctx) error
	UpdateSize(c *fiber.Ctx) error
}

func NewSizeHandler() SizeHandlerInterface {
	return &sizeHandler{}
}

var SizeHandler = NewSizeHandler()

func (h *sizeHandler) GetAllSizes(c *fiber.Ctx) error { // Retrieve pagination parameters from context.
	pageNumberInterface := c.Locals("pageNumber")
	pageSizeInterface := c.Locals("pageSize")

	pageNumber, ok := pageNumberInterface.(int)
	if !ok {
		log.Error("Failed to retrieve page number from context")
		pageNumber = 1
	}

	pageSize, ok := pageSizeInterface.(int)
	if !ok {
		log.Error("Failed to retrieve page size from context")
		pageSize = 100
	}

	users, err := service.SizeService.GetAllSizes(pageNumber, pageSize)
	if err != nil {
		log.Error("Failed to fetch paginated sizes", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch sizes", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *sizeHandler) CreateSize(c *fiber.Ctx) error {
	reqInterface := c.Locals("validatedBody")

	body, ok := reqInterface.(dto.CreateSizeRequest)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	size, err := service.SizeService.CreateSize(&body)
	if err != nil {
		log.Error("Failed to create size", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create size", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(size)
}

func (h *sizeHandler) UpdateSize(c *fiber.Ctx) error {
	reqDto := c.Locals("validatedBody")

	body, ok := reqDto.(dto.UpdateSizeRequest)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	size, err := service.SizeService.UpdateSize(&body)
	if err != nil {
		log.Error("Failed to create size", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create size", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(size)
}

func (h *sizeHandler) DeleteSize(c *fiber.Ctx) error {

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

	err = service.SizeService.DeleteSize(sizeId)
	if err != nil {
		log.Error("Failed to remove size", zap.Int("sizeId", sizeId), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to remove user", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": sizeId})
}
