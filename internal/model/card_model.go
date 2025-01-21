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
	NodeId              int              `json:"nodeId"`
	Title               string           `json:"title"`
	NodeDescription     *string          `json:"nodeDescription"`
	CreatedAt           string           `json:"createdAt"`
	UpdatedAt           string           `json:"updatedAt"`
	RemovedAt           *string          `json:"removedAt"`
	Images              []string         `db:"images" json:"images"`
	NodeType            string           `json:"nodeType"`
	NodeTypeDescription *string          `json:"nodeTypeDescription"`
	Characteristics     []Characteristic `json:"characteristics"`
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

	m := make(map[int]CardResponse)

	for _, row := range *rows {
		nodeId := row.NodeId
		value, exist := m[nodeId]

		addParam := Characteristic{
			Title:                     row.Characteristic,
			Value:                     row.CharacteristicValue,
			CharacteristicDescription: row.CharacteristicDescription,
		}

		if row.AdditionalParams != nil {
			var parsedField map[string]interface{}
			if err := json.Unmarshal(*row.AdditionalParams, &parsedField); err != nil {
				return nil, fmt.Errorf("failed to parse additionalParams for NodeId %d: %w", nodeId, err)
			}
			addParam.AdditionalParams = &parsedField
		}

		if exist {
			value.Characteristics = append(value.Characteristics, addParam)
			m[nodeId] = value
		} else {
			m[nodeId] = CardResponse{
				NodeId:              nodeId,
				Title:               row.Title,
				NodeDescription:     row.NodeDescription,
				CreatedAt:           row.CreatedAt,
				UpdatedAt:           row.UpdatedAt,
				RemovedAt:           row.RemovedAt,
				Images:              row.Images,
				NodeType:            row.NodeType,
				NodeTypeDescription: row.NodeTypeDescription,
				Characteristics:     []Characteristic{addParam},
			}
		}
	}

	result := make([]CardResponse, 0, len(m))
	for _, cardResponse := range m {
		result = append(result, cardResponse)
	}

	return result, nil
}
