package service

import (
	"go.uber.org/zap"
	"shop/internal/api/dto"
	"shop/internal/model"
	"shop/internal/repository"
	"shop/pkg/log"
	"shop/pkg/utils"
)

type sizeService struct{}

type SizeServiceInterface interface {
	GetAllSizes(pageNumber, pageSize int) (*model.Paginate[model.SizeRow], error)
	CreateSize(size *dto.CreateSizeRequest) (*model.SizeRow, error)
	UpdateSize(size *dto.UpdateSizeRequest) (*model.SizeRow, error)
	DeleteSize(id int) error
}

func NewSizeService() SizeServiceInterface {
	return &sizeService{}
}

var SizeService = NewSizeService()

func (s *sizeService) GetAllSizes(pageNumber, pageSize int) (*model.Paginate[model.SizeRow], error) {
	sizes, totalCount, err := repository.SizeRepo.GetAllSizes(pageNumber, pageSize)
	if err != nil {
		log.Error("Failed to fetch sizes", zap.Error(err))
		return nil, err
	}

	// Формируем структуру Paginate с типом SizeRow
	result := &model.Paginate[model.SizeRow]{
		PageNumber:     pageNumber,
		RowTotalCount:  totalCount,
		TotalPageCount: utils.CalculateTotalPages(totalCount, pageSize),
		PageSize:       pageSize,
		Items:          sizes,
	}

	return result, nil
}

func (s *sizeService) CreateSize(size *dto.CreateSizeRequest) (*model.SizeRow, error) {

	createdID, err := repository.SizeRepo.CreateSize(size)
	if err != nil {
		log.Error("Failed to create size", zap.Error(err))
		return nil, err
	}

	sizeRow, err := repository.SizeRepo.GetSizeById(createdID)
	if err != nil {
		return nil, err
	}

	return sizeRow, nil
}

func (s *sizeService) UpdateSize(dto *dto.UpdateSizeRequest) (*model.SizeRow, error) {
	_, err := repository.SizeRepo.GetSizeById(dto.ID)
	if err != nil {
		return nil, err
	}

	row := model.SizeRow{
		ID:          dto.ID,
		Title:       dto.Title,
		Description: dto.Description,
	}

	if err := repository.SizeRepo.UpdateSize(&row); err != nil {
		log.Error("Failed to update size", zap.Error(err))
		return nil, err
	}
	return &row, nil
}

func (s *sizeService) DeleteSize(id int) error {
	if err := repository.SizeRepo.DeleteSizeById(id); err != nil {
		log.Error("Failed to delete size", zap.Error(err))
		return err
	}
	return nil
}
