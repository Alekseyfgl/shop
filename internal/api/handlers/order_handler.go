package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"shop/internal/api/dto"
	"shop/pkg/http_error"
	"shop/pkg/log"
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

	_, ok := reqInterface.(dto.OrderDTO)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	orderId := uuid.New().String()
	// здесь можно вызвать сервис создания заказа и т.д.

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"orderId": orderId,
	})
}
