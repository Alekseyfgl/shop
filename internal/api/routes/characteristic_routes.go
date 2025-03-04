package routes

import (
	"github.com/gofiber/fiber/v2"
	"shop/internal/api/handlers"
	"shop/internal/api/middlewares/validator/dto_validator"
)

func RegisterCharacteristicRoutes(app fiber.Router) {
	app.Get("/characteristics",
		dto_validator.ValidatePaginationMiddleware(),
		handlers.CharacteristicHandler.GetAllCharacteristics,
	)
	app.Get("/characteristics/filters",
		dto_validator.ValidateNodeTypeIdMiddleware(),
		handlers.CharacteristicHandler.GetCharForFilters,
	)
	app.Post("/characteristics",
		dto_validator.ValidateCreateCharacteristicMiddleware(),
		handlers.CharacteristicHandler.CreateCharacteristic,
	)

	app.Put("/characteristics",
		dto_validator.ValidateUpdateCharacteristicMiddleware(),
		handlers.CharacteristicHandler.UpdateCharacteristic,
	)
	app.Delete("/characteristics/:id",
		dto_validator.ValidateIdMiddleware(),
		handlers.CharacteristicHandler.DeleteCharacteristic,
	)
}
