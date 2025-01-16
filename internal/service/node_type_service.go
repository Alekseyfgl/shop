package service

import (
	"go.uber.org/zap"
	"shop/internal/api/dto"
	"shop/internal/model"
	"shop/internal/repository"
	"shop/pkg/log"
	"shop/pkg/utils"
)

type nodeTypeService struct{}

type NodeTypeServiceInterface interface {
	GetAllNodeType(pageNumber, pageSize int) (*model.Paginate[model.NodeTypeRow], error)
	CreateNodeType(size *dto.CreateNodeTypeRequest) (*model.NodeTypeRow, error)
	UpdateNodeType(size *dto.UpdateNodeTypeRequest) (*model.NodeTypeRow, error)
	DeleteNodeType(id int) error
}

func NewNodeTypeService() NodeTypeServiceInterface {
	return &nodeTypeService{}
}

var NodeTypeService = NewNodeTypeService()

func (s *nodeTypeService) GetAllNodeType(pageNumber, pageSize int) (*model.Paginate[model.NodeTypeRow], error) {
	nodeTypes, totalCount, err := repository.NodeTypeRepo.GetAllNodeTypes(pageNumber, pageSize)
	if err != nil {
		log.Error("Failed to fetch node_type", zap.Error(err))
		return nil, err
	}

	// Формируем структуру Paginate с типом SizeRow
	result := &model.Paginate[model.NodeTypeRow]{
		PageNumber:     pageNumber,
		RowTotalCount:  totalCount,
		TotalPageCount: utils.CalculateTotalPages(totalCount, pageSize),
		PageSize:       pageSize,
		Items:          nodeTypes,
	}

	return result, nil
}

func (s *nodeTypeService) CreateNodeType(dto *dto.CreateNodeTypeRequest) (*model.NodeTypeRow, error) {

	createdID, err := repository.NodeTypeRepo.CreateNodeType(dto)
	if err != nil {
		log.Error("Failed to create node_type", zap.Error(err))
		return nil, err
	}

	nodeType, err := repository.NodeTypeRepo.GetNodeTypeById(createdID)
	if err != nil {
		return nil, err
	}

	return nodeType, nil
}

func (s *nodeTypeService) UpdateNodeType(dto *dto.UpdateNodeTypeRequest) (*model.NodeTypeRow, error) {
	_, err := repository.NodeTypeRepo.GetNodeTypeById(dto.ID)
	if err != nil {
		return nil, err
	}

	row := model.NodeTypeRow{
		ID:          dto.ID,
		Type:        dto.Type,
		Description: dto.Description,
	}

	if err := repository.NodeTypeRepo.UpdateNodeType(&row); err != nil {
		log.Error("Failed to update node_type", zap.Error(err))
		return nil, err
	}
	return &row, nil
}

func (s *nodeTypeService) DeleteNodeType(id int) error {
	if err := repository.NodeTypeRepo.DeleteNodeTypeById(id); err != nil {
		log.Error("Failed to delete characteristics", zap.Error(err))
		return err
	}
	return nil
}
