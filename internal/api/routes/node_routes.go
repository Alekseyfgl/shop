package routes

import (
	"github.com/gofiber/fiber/v2"
	"shop/internal/api/handlers"
	"shop/internal/api/middlewares/validator/dto_validator"
)

func RegisterNodeRoutes(app fiber.Router) {
	app.Get("/nodes",
		dto_validator.ValidatePaginationMiddleware(),
		handlers.NodeHandler.GetAllNode,
	)
	app.Post("/nodes/create",
		dto_validator.ValidateCreateNodeMiddleware(),
		handlers.NodeHandler.CreateNode,
	)

	app.Put("/nodes",
		dto_validator.ValidateUpdateNodeMiddleware(),
		handlers.NodeHandler.UpdateNode,
	)
	app.Delete("/nodes/:id",
		dto_validator.ValidateIdMiddleware(),
		handlers.NodeHandler.DeleteNode,
	)
}
