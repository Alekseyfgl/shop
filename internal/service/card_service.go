package service

import (
	"errors"
	"go.uber.org/zap"
	"shop/internal/api/dto"
	"shop/internal/model"
	"shop/internal/repository"
	"shop/pkg/log"
	"shop/pkg/utils"
)

type cardService struct{}

type CardServiceInterface interface {
	GetCardById(id int) (*model.CardResponse, error)
	GetAllCards(pageNumber, pageSize int, filters *[]model.CardFilter) (*model.Paginate[model.CardResponse], error)
	CreateCard(dto *dto.CreateCardDTO) (*model.CardResponse, error)
	GetCardsByVector(dto *dto.GetCardsByVectorDTO) (*[]model.CardResponse, error)
}

func NewCardService() CardServiceInterface {
	return &cardService{}
}

var CardService = NewCardService()

func (s *cardService) GetCardById(id int) (*model.CardResponse, error) {
	card, err := repository.CardRepo.GetCardById(id)
	if err != nil {
		return nil, err
	}

	result, err := model.MapperCardResponse(card)
	if err != nil {
		return nil, err
	}

	el := &result[0]
	return el, nil
}

func (s *cardService) GetAllCards(pageNumber, pageSize int, filters *[]model.CardFilter) (*model.Paginate[model.CardResponse], error) {
	cards, totalCount, err := repository.CardRepo.GetAllCards(pageNumber, pageSize, filters)
	if err != nil {
		log.Error("Failed to fetch cards", zap.Error(err))
		return nil, err
	}

	mappedCards, err := model.MapperCardResponse(cards)
	if err != nil {
		return nil, err
	}

	result := &model.Paginate[model.CardResponse]{
		PageNumber:     pageNumber,
		RowTotalCount:  totalCount,
		TotalPageCount: utils.CalculateTotalPages(totalCount, pageSize),
		PageSize:       pageSize,
		Items:          mappedCards,
	}

	return result, nil
}

func (s *cardService) GetCardsByVector(dto *dto.GetCardsByVectorDTO) (*[]model.CardResponse, error) {
	cards, err := repository.CardRepo.FindByVectorSearch(dto.Text, dto.Limit)
	if err != nil {
		log.Error("Failed to fetch cards", zap.Error(err))
		return nil, err
	}

	mappedCards, err := model.MapperCardResponse(cards)
	if err != nil {
		return nil, err
	}

	return &mappedCards, nil
}

func (s *cardService) CreateCard(dto *dto.CreateCardDTO) (*model.CardResponse, error) {
	newID, err := repository.CardRepo.CreateCard(dto)
	if err != nil {
		log.Error("Failed to fetch card, after creating", zap.Error(err))
		return nil, errors.New("failed to create card")
	}

	newCard, err := s.GetCardById(newID)
	if err != nil {
		log.Error("Failed to fetch card, after creating", zap.Error(err))
		return nil, err
	}
	return newCard, nil
}
