package routes

import (
	"github.com/gofiber/fiber/v2"
	"shop/internal/api/handlers"
	"shop/internal/api/middlewares/validator/dto_validator"
)

func RegisterCharDefaultValueRoutes(app fiber.Router) {
	app.Get("/selectors",
		dto_validator.ValidatePaginationMiddleware(),
		handlers.CharDefaultValueHandler.GetAllDefValue,
	)

	app.Get("/selectors/:id",
		dto_validator.ValidateIdMiddleware(),
		handlers.CharDefaultValueHandler.GetDefValueById,
	)
	app.Post("/selectors",
		dto_validator.ValidateCreateCharDefaultValueMiddleware(),
		handlers.CharDefaultValueHandler.CreateDefValue,
	)

	app.Put("/selectors",
		dto_validator.ValidateUpdateCharDefValueMiddleware(),
		handlers.CharDefaultValueHandler.UpdateDefValue,
	)
	app.Delete("/selectors/:id",
		dto_validator.ValidateIdMiddleware(),
		handlers.CharDefaultValueHandler.DeleteDefValue,
	)
}
