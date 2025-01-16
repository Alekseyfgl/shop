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

type nodeTypeHandler struct{}

type NodeTypeHandlerInterface interface {
	GetAllNodeType(c *fiber.Ctx) error
	CreateNodeType(c *fiber.Ctx) error
	DeleteNodeType(c *fiber.Ctx) error
	UpdateNodeType(c *fiber.Ctx) error
}

func NewNodeTypeHandler() NodeTypeHandlerInterface {
	return &nodeTypeHandler{}
}

var NodeTypeHandler = NewNodeTypeHandler()

func (h *nodeTypeHandler) GetAllNodeType(c *fiber.Ctx) error {
	pageNumberInterface := c.Locals("pageNumber")
	pageSizeInterface := c.Locals("pageSize")

	pageNumber, ok := pageNumberInterface.(int)
	if !ok {
		log.Error("Failed to retrieve page number from context")
		pageNumber = 1
	}

	pageSize, ok := pageSizeInterface.(int)
	if !ok {
		log.Error("Failed to retrieve page node_type from context")
		pageSize = 100
	}

	users, err := service.NodeTypeService.GetAllNodeType(pageNumber, pageSize)
	if err != nil {
		log.Error("Failed to fetch paginated node_type", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch node_type", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *nodeTypeHandler) CreateNodeType(c *fiber.Ctx) error {
	reqInterface := c.Locals("validatedBody")

	body, ok := reqInterface.(dto.CreateNodeTypeRequest)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	size, err := service.NodeTypeService.CreateNodeType(&body)
	if err != nil {
		log.Error("Failed to create node_type", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create node_type", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(size)
}

func (h *nodeTypeHandler) UpdateNodeType(c *fiber.Ctx) error {
	reqDto := c.Locals("validatedBody")

	body, ok := reqDto.(dto.UpdateNodeTypeRequest)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	size, err := service.NodeTypeService.UpdateNodeType(&body)
	if err != nil {
		log.Error("Failed to create node_type", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create node_type", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(size)
}

func (h *nodeTypeHandler) DeleteNodeType(c *fiber.Ctx) error {

	// Retrieve article ID from context.
	id := c.Locals("Id")
	nodeTypeIdStr, ok := id.(string)
	if !ok || nodeTypeIdStr == "0" {
		log.Error("ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "ID is required", nil).Send(c)
	}

	nodeTypeId, err := utils.StringToInt(nodeTypeIdStr)
	if err != nil {
		return http_error.NewHTTPError(fiber.StatusInternalServerError, err.Error(), nil).Send(c)
	}

	err = service.NodeTypeService.DeleteNodeType(nodeTypeId)
	if err != nil {
		log.Error("Failed to remove node_type", zap.Int("nodeTypeId", nodeTypeId), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to remove user", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": nodeTypeId})
}
