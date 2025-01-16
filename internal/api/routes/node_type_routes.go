package routes

import (
	"github.com/gofiber/fiber/v2"
	"shop/internal/api/handlers"
	"shop/internal/api/middlewares/validator/dto_validator"
)

func RegisterNodeTypeRoutes(app *fiber.App) {
	app.Get("/node-types",
		dto_validator.ValidatePaginationMiddleware(),
		handlers.NodeTypeHandler.GetAllNodeType,
	)
	app.Post("/node-types",
		dto_validator.ValidateCreateNodeTypeMiddleware(),
		handlers.NodeTypeHandler.CreateNodeType,
	)

	app.Put("/node-types",
		dto_validator.ValidateUpdateNodeTypeMiddleware(),
		handlers.NodeTypeHandler.UpdateNodeType,
	)
	app.Delete("/node-types/:id",
		dto_validator.ValidateIdMiddleware(),
		handlers.NodeTypeHandler.DeleteNodeType,
	)
}
