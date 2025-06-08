package handlers

import (
	"shop/internal/api/dto"
	"shop/pkg/http_error"
	"shop/pkg/log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type orderHandler struct{}

type OrderHandlerInterface interface {
	CreateOrder(c *fiber.Ctx) error
}

func NewOrderHandler() OrderHandlerInterface {
	return &orderHandler{}
}

var OrderHandler = NewOrderHandler()

func (h *orderHandler) CreateOrder(c *fiber.Ctx) error {
	reqInterface := c.Locals("validatedBody")

	_, ok := reqInterface.([]dto.OrderDTO)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	// Генерируем ID для каждого заказа

	orderId := uuid.New().String()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"orderId": orderId,
	})
}
