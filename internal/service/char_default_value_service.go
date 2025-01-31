package service

import (
	"go.uber.org/zap"
	"shop/internal/api/dto"
	"shop/internal/model"
	"shop/internal/repository"
	"shop/pkg/log"
	"shop/pkg/utils"
)

type charDefaultValueService struct{}

type CharDefaultValueServiceInterface interface {
	GetAllDefValue(pageNumber, pageSize int) (*model.Paginate[model.CharDefaultValue], error)
	GetDefValueById(id int) (*[]model.CharDefaultValue, error)
	CreateDefValue(size *dto.CreateCharDefValueRequest) (*model.CharDefaultValueRow, error)
	UpdateDefValue(size *dto.UpdateCharDefValueRequest) (*model.CharDefaultValueRow, error)
	DeleteDefValueById(id int) error
}

func NewCharDefaultValueService() CharDefaultValueServiceInterface {
	return &charDefaultValueService{}
}

var CharDefaultValueService = NewCharDefaultValueService()

func (s *charDefaultValueService) GetAllDefValue(pageNumber, pageSize int) (*model.Paginate[model.CharDefaultValue], error) {
	defValues, totalCount, err := repository.CharDefaultValueRepo.GetAllDefaultValues(pageNumber, pageSize)
	if err != nil {
		log.Error("Failed to fetch char_default_value", zap.Error(err))
		return nil, err
	}

	// Формируем структуру Paginate с типом SizeRow
	result := &model.Paginate[model.CharDefaultValue]{
		PageNumber:     pageNumber,
		RowTotalCount:  totalCount,
		TotalPageCount: utils.CalculateTotalPages(totalCount, pageSize),
		PageSize:       pageSize,
		Items:          defValues,
	}

	return result, nil
}

func (s *charDefaultValueService) GetDefValueById(id int) (*[]model.CharDefaultValue, error) {
	defValues, err := repository.CharDefaultValueRepo.GetFullDefaultValueById(id)
	if err != nil {
		log.Error("Failed to fetch char_default_value", zap.Error(err))
		return nil, err
	}
	return defValues, nil
}

func (s *charDefaultValueService) CreateDefValue(dto *dto.CreateCharDefValueRequest) (*model.CharDefaultValueRow, error) {

	createdID, err := repository.CharDefaultValueRepo.CreateDefaultValue(dto)
	if err != nil {
		log.Error("Failed to create char_default_value", zap.Error(err))
		return nil, err
	}

	nodeType, err := repository.CharDefaultValueRepo.GetDefaultValueById(createdID)
	if err != nil {
		return nil, err
	}

	return nodeType, nil
}

func (s *charDefaultValueService) UpdateDefValue(dto *dto.UpdateCharDefValueRequest) (*model.CharDefaultValueRow, error) {
	_, err := repository.CharDefaultValueRepo.GetDefaultValueById(dto.ID)
	if err != nil {
		return nil, err
	}

	if err := repository.CharDefaultValueRepo.UpdateDefaultValue(dto); err != nil {
		log.Error("Failed to update char_default_value", zap.Error(err))
		return nil, err
	}

	updatedRow, err := repository.CharDefaultValueRepo.GetDefaultValueById(dto.ID)
	if err != nil {
		return nil, err
	}

	return updatedRow, nil
}

func (s *charDefaultValueService) DeleteDefValueById(id int) error {
	if err := repository.CharDefaultValueRepo.DeleteDefaultValueById(id); err != nil {
		log.Error("Failed to delete char_default_value", zap.Error(err))
		return err
	}
	return nil
}
