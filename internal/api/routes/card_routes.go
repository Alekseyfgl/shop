package routes

import (
	"github.com/gofiber/fiber/v2"
	"shop/internal/api/handlers"
	"shop/internal/api/middlewares/validator/dto_validator"
)

func RegisterCardRoutes(app fiber.Router) {
	app.Get("/cards/:id",
		dto_validator.ValidateIdMiddleware(),
		handlers.CardHandler.GetCardById,
	)
	app.Get("/cards",
		dto_validator.ValidateNodeTypeIdMiddleware(),
		dto_validator.ValidatePaginationMiddleware(),
		handlers.CardHandler.GetAllCards,
	)
	app.Post("/cards",
		dto_validator.ValidateCreateCardMiddleware(),
		handlers.CardHandler.CreateCard,
	)
	app.Post("/cards/search",
		dto_validator.ValidateGetCardsByVectorMiddleware(),
		handlers.CardHandler.GetCardsByVector,
	)
}
