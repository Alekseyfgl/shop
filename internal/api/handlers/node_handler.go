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

type nodeHandler struct{}

type NodeHandlerInterface interface {
	GetAllNode(c *fiber.Ctx) error
	CreateNode(c *fiber.Ctx) error
	DeleteNode(c *fiber.Ctx) error
	UpdateNode(c *fiber.Ctx) error
}

func NewNodeHandler() NodeHandlerInterface {
	return &nodeHandler{}
}

var NodeHandler = NewNodeHandler()

func (h *nodeHandler) GetAllNode(c *fiber.Ctx) error {
	pageNumberInterface := c.Locals("pageNumber")
	pageSizeInterface := c.Locals("pageSize")

	pageNumber, ok := pageNumberInterface.(int)
	if !ok {
		log.Error("Failed to retrieve page number from context")
		pageNumber = 1
	}

	pageSize, ok := pageSizeInterface.(int)
	if !ok {
		log.Error("Failed to retrieve page node from context")
		pageSize = 100
	}

	users, err := service.NodeService.GetAllNode(pageNumber, pageSize)
	if err != nil {
		log.Error("Failed to fetch paginated node", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch node", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

func (h *nodeHandler) CreateNode(c *fiber.Ctx) error {
	reqInterface := c.Locals("validatedBody")

	body, ok := reqInterface.(dto.CreateNodeRequest)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	size, err := service.NodeService.CreateNode(&body)
	if err != nil {
		log.Error("Failed to create node", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create node", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(size)
}

func (h *nodeHandler) UpdateNode(c *fiber.Ctx) error {
	reqDto := c.Locals("validatedBody")

	body, ok := reqDto.(dto.UpdateNodeRequest)
	if !ok {
		log.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	size, err := service.NodeService.UpdateNode(&body)
	if err != nil {
		log.Error("Failed to create node", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create node", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(size)
}

func (h *nodeHandler) DeleteNode(c *fiber.Ctx) error {

	// Retrieve article ID from context.
	id := c.Locals("Id")
	nodeIdIdStr, ok := id.(string)
	if !ok || nodeIdIdStr == "0" {
		log.Error("ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "ID is required", nil).Send(c)
	}

	nodeId, err := utils.StringToInt(nodeIdIdStr)
	if err != nil {
		return http_error.NewHTTPError(fiber.StatusInternalServerError, err.Error(), nil).Send(c)
	}

	err = service.NodeService.DeleteNode(nodeId)
	if err != nil {
		log.Error("Failed to remove node_type", zap.Int("nodeId", nodeId), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to remove user", nil).Send(c)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": nodeId})
}
