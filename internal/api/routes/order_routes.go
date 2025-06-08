package routes

import (
	"github.com/gofiber/fiber/v2"
	"shop/internal/api/handlers"
	"shop/internal/api/middlewares/validator/dto_validator"
)

func RegisterOrderRoutes(app fiber.Router) {

	app.Post("/orders",
		dto_validator.ValidateCreateOrderMiddleware(),
		handlers.OrderHandler.CreateOrder,
	)

}
