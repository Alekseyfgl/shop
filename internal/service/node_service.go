package service

import (
	"go.uber.org/zap"
	"shop/internal/api/dto"
	"shop/internal/model"
	"shop/internal/repository"
	"shop/pkg/log"
	"shop/pkg/utils"
)

type nodeService struct{}

type NodeServiceInterface interface {
	GetAllNode(pageNumber, pageSize int) (*model.Paginate[model.NodeRow], error)
	CreateNode(size *dto.CreateNodeRequest) (*model.NodeRow, error)
	UpdateNode(size *dto.UpdateNodeRequest) (*model.NodeRow, error)
	DeleteNode(id int) error
}

func NewNodeService() NodeServiceInterface {
	return &nodeService{}
}

var NodeService = NewNodeService()

func (s *nodeService) GetAllNode(pageNumber, pageSize int) (*model.Paginate[model.NodeRow], error) {
	nodeTypes, totalCount, err := repository.NodeRepo.GetAllNodes(pageNumber, pageSize)
	if err != nil {
		log.Error("Failed to fetch node", zap.Error(err))
		return nil, err
	}

	// Формируем структуру Paginate с типом SizeRow
	result := &model.Paginate[model.NodeRow]{
		PageNumber:     pageNumber,
		RowTotalCount:  totalCount,
		TotalPageCount: utils.CalculateTotalPages(totalCount, pageSize),
		PageSize:       pageSize,
		Items:          nodeTypes,
	}

	return result, nil
}

func (s *nodeService) CreateNode(dto *dto.CreateNodeRequest) (*model.NodeRow, error) {
	log.Info("START func (s *nodeService) CreateNode(dto *dto.CreateNodeRequest) (*model.NodeRow, error)")
	_, err := repository.NodeTypeRepo.GetNodeTypeById(dto.NodeTypeId)
	if err != nil {
		log.Error("NodeTypeId not found", zap.Error(err))
		return nil, err
	}
	log.Info("repository.NodeRepo.CreateNode(dto)")
	createdID, err := repository.NodeRepo.CreateNode(dto)
	if err != nil {
		log.Error("Failed to create node", zap.Error(err))
		return nil, err
	}
	log.Info("repository.NodeRepo.GetNodeById(createdID)")
	nodeType, err := repository.NodeRepo.GetNodeById(createdID)
	if err != nil {
		return nil, err
	}

	return nodeType, nil
}

func (s *nodeService) UpdateNode(dto *dto.UpdateNodeRequest) (*model.NodeRow, error) {
	_, err := repository.NodeRepo.GetNodeById(dto.ID)
	if err != nil {
		return nil, err
	}

	_, err = repository.NodeTypeRepo.GetNodeTypeById(dto.NodeTypeId)
	if err != nil {
		log.Error("NodeTypeId not found", zap.Error(err))
		return nil, err
	}

	if err := repository.NodeRepo.UpdateNodes(dto); err != nil {
		log.Error("Failed to update node", zap.Error(err))
		return nil, err
	}

	updatedNode, err := repository.NodeRepo.GetNodeById(dto.ID)
	if err != nil {
		log.Error("Failed to update node", zap.Error(err))
		return nil, err
	}

	return updatedNode, nil
}

func (s *nodeService) DeleteNode(id int) error {
	if err := repository.NodeRepo.DeleteNodeById(id); err != nil {
		log.Error("Failed to delete node", zap.Error(err))
		return err
	}
	return nil
}
