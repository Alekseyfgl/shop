package routes

import (
	"github.com/gofiber/fiber/v2"
	"shop/internal/api/handlers"
	"shop/internal/api/middlewares/validator/dto_validator"
)

func RegisterSizeRoutes(app *fiber.App) {
	app.Get("/sizes",
		dto_validator.ValidatePaginationMiddleware(),
		handlers.SizeHandler.GetAllSizes,
	)
	app.Post("/sizes",
		dto_validator.ValidateCreateSizeMiddleware(),
		handlers.SizeHandler.CreateSize,
	)

	app.Put("/sizes",
		dto_validator.ValidateUpdateSizeMiddleware(),
		handlers.SizeHandler.UpdateSize,
	)
	app.Delete("/sizes/:id",
		dto_validator.ValidateIdMiddleware(),
		handlers.SizeHandler.DeleteSize,
	)
}
