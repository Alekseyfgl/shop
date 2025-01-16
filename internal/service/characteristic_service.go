package service

import (
	"go.uber.org/zap"
	"shop/internal/api/dto"
	"shop/internal/model"
	"shop/internal/repository"
	"shop/pkg/log"
	"shop/pkg/utils"
)

type characteristicService struct{}

type CharacteristicServiceInterface interface {
	GetAllCharacteristic(pageNumber, pageSize int) (*model.Paginate[model.CharacteristicRow], error)
	GetCharForFilters() (*[]model.CharFilterResponse, error)
	CreateCharacteristic(size *dto.CreateCharacteristicRequest) (*model.CharacteristicRow, error)
	UpdateCharacteristic(size *dto.UpdateCharacteristicRequest) (*model.CharacteristicRow, error)
	DeleteCharacteristic(id int) error
}

func NewCharacteristicService() CharacteristicServiceInterface {
	return &characteristicService{}
}

var CharacteristicService = NewCharacteristicService()

func (s *characteristicService) GetAllCharacteristic(pageNumber, pageSize int) (*model.Paginate[model.CharacteristicRow], error) {
	characteristics, totalCount, err := repository.CharacteristicRepo.GetAllCharacteristics(pageNumber, pageSize)
	if err != nil {
		log.Error("Failed to fetch characteristics", zap.Error(err))
		return nil, err
	}

	// Формируем структуру Paginate с типом SizeRow
	result := &model.Paginate[model.CharacteristicRow]{
		PageNumber:     pageNumber,
		RowTotalCount:  totalCount,
		TotalPageCount: utils.CalculateTotalPages(totalCount, pageSize),
		PageSize:       pageSize,
		Items:          characteristics,
	}

	return result, nil
}

func (s *characteristicService) CreateCharacteristic(dto *dto.CreateCharacteristicRequest) (*model.CharacteristicRow, error) {

	createdID, err := repository.CharacteristicRepo.CreateCharacteristics(dto)
	if err != nil {
		log.Error("Failed to create characteristic", zap.Error(err))
		return nil, err
	}

	characteristic, err := repository.CharacteristicRepo.GetCharacteristicsById(createdID)
	if err != nil {
		return nil, err
	}

	return characteristic, nil
}

func (s *characteristicService) UpdateCharacteristic(dto *dto.UpdateCharacteristicRequest) (*model.CharacteristicRow, error) {
	_, err := repository.CharacteristicRepo.GetCharacteristicsById(dto.ID)
	if err != nil {
		return nil, err
	}

	row := model.CharacteristicRow{
		ID:          dto.ID,
		Title:       dto.Title,
		Description: dto.Description,
	}

	if err := repository.CharacteristicRepo.UpdateCharacteristics(&row); err != nil {
		log.Error("Failed to update size", zap.Error(err))
		return nil, err
	}
	return &row, nil
}

func (s *characteristicService) DeleteCharacteristic(id int) error {
	if err := repository.CharacteristicRepo.DeleteCharacteristicsById(id); err != nil {
		log.Error("Failed to delete characteristics", zap.Error(err))
		return err
	}
	return nil
}

func (s *characteristicService) GetCharForFilters() (*[]model.CharFilterResponse, error) {
	// Получаем данные из репозитория
	filters, err := repository.CharacteristicRepo.GetCharFilters()
	if err != nil {
		return nil, err
	}

	// Создаем временную мапу для группировки
	grouped := make(map[int]*model.CharFilterResponse)

	// Группируем данные
	for _, row := range *filters {
		if _, exists := grouped[row.CharacteristicId]; !exists {
			grouped[row.CharacteristicId] = &model.CharFilterResponse{
				CharacteristicId: row.CharacteristicId,
				Title:            row.Title,
				Values:           []string{},
			}
		}

		// Добавляем значение к существующей записи, если значение не null
		if row.Value != nil {
			grouped[row.CharacteristicId].Values = append(grouped[row.CharacteristicId].Values, *row.Value)
		}
	}

	// Преобразуем мапу в массив
	result := make([]model.CharFilterResponse, 0, len(grouped))
	for _, value := range grouped {
		result = append(result, *value)
	}

	return &result, nil
}
