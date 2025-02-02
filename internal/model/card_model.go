package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

type CardRow struct {
	NodeId                    int              `db:"nodeId" json:"nodeId"`
	Title                     string           `db:"title" json:"title"`
	NodeDescription           *string          `db:"nodeDescription" json:"nodeDescription"`
	CreatedAt                 string           `db:"createdAt" json:"createdAt"`
	UpdatedAt                 string           `db:"updatedAt" json:"updatedAt"`
	RemovedAt                 *string          `db:"removedAt" json:"removedAt"`
	Images                    []string         `db:"images" json:"images"`
	NodeType                  string           `db:"nodeType" json:"nodeType"`
	NodeTypeDescription       *string          `db:"nodeTypeDescription" json:"nodeTypeDescription"`
	Characteristic            string           `db:"characteristic" json:"characteristic"`
	CharacteristicValue       string           `db:"characteristicValue" json:"characteristicValue"`
	AdditionalParams          *json.RawMessage `db:"additionalParams" json:"additionalParams"`
	CharacteristicDescription *string          `db:"characteristicDescription" json:"characteristicDescription"`
}
type CardResponse struct {
	NodeId              int                `json:"nodeId"`
	Title               string             `json:"title"`
	NodeDescription     *string            `json:"nodeDescription"`
	CreatedAt           string             `json:"createdAt"`
	UpdatedAt           string             `json:"updatedAt"`
	RemovedAt           *string            `json:"removedAt"`
	Images              []string           `db:"images" json:"images"`
	NodeType            string             `json:"nodeType"`
	NodeTypeDescription *string            `json:"nodeTypeDescription"`
	Characteristics     [][]Characteristic `json:"characteristics"`
}

type Characteristic struct {
	Title                     string                  `json:"title"`
	Value                     string                  `json:"value"`
	AdditionalParams          *map[string]interface{} `json:"additionalParams"`
	CharacteristicDescription *string                 `json:"description"`
}

type CardFilter struct {
	Key    string
	Values string
}

func MapperCardResponse(rows *[]CardRow) ([]CardResponse, error) {
	if len(*rows) == 0 {
		return nil, errors.New("card not found")
	}

	// Вспомогательная структура для группировки характеристик по title для каждой ноды.
	type cardGroup struct {
		card   CardResponse
		groups map[string][]Characteristic
	}

	// Используем карту, где ключ – NodeId, а значение – указатель на cardGroup.
	m := make(map[int]*cardGroup)

	for _, row := range *rows {
		nodeId := row.NodeId

		// Формируем объект характеристики из строки.
		addParam := Characteristic{
			Title:                     row.Characteristic,
			Value:                     row.CharacteristicValue,
			CharacteristicDescription: row.CharacteristicDescription,
		}

		// Если есть дополнительные параметры – парсим их.
		if row.AdditionalParams != nil {
			var parsedField map[string]interface{}
			if err := json.Unmarshal(*row.AdditionalParams, &parsedField); err != nil {
				return nil, fmt.Errorf("failed to parse additionalParams for NodeId %d: %w", nodeId, err)
			}
			addParam.AdditionalParams = &parsedField
		}

		// Если для данной ноды уже создана карточка – просто добавляем характеристику в нужную группу.
		if cg, exist := m[nodeId]; exist {
			cg.groups[addParam.Title] = append(cg.groups[addParam.Title], addParam)
		} else {
			// Иначе создаём новую карточку и инициализируем группу характеристик.
			newCard := CardResponse{
				NodeId:              row.NodeId,
				Title:               row.Title,
				NodeDescription:     row.NodeDescription,
				CreatedAt:           row.CreatedAt,
				UpdatedAt:           row.UpdatedAt,
				RemovedAt:           row.RemovedAt,
				Images:              row.Images,
				NodeType:            row.NodeType,
				NodeTypeDescription: row.NodeTypeDescription,
				// Поле Characteristics пока оставляем nil – заполнится после группировки.
			}
			groups := make(map[string][]Characteristic)
			groups[addParam.Title] = []Characteristic{addParam}
			m[nodeId] = &cardGroup{
				card:   newCard,
				groups: groups,
			}
		}
	}

	// Формируем итоговый срез карточек, где характеристики сгруппированы в подмассивы.
	result := make([]CardResponse, 0, len(m))
	for _, cg := range m {
		// Преобразуем карту групп в срез срезов.
		var groupedCharacteristics [][]Characteristic
		for _, chars := range cg.groups {
			groupedCharacteristics = append(groupedCharacteristics, chars)
		}
		cg.card.Characteristics = groupedCharacteristics
		result = append(result, cg.card)
	}

	return result, nil
}
