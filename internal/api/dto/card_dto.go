package dto

import "encoding/json"

type CreateCardDTO struct {
	Title           string    `json:"title" validate:"required,min=1"`
	NodeDescription *string   `json:"nodeDescription" validate:"omitempty,min=3,max=1000"`
	NodeTypeId      int       `json:"nodeTypeId" validate:"required,number"`
	Images          []string  `json:"images" validate:"required,min=1,dive"`
	Characteristics []CharDTO `json:"characteristics" validate:"required,min=1,dive"`
}

type CharDTO struct {
	Id               int              `json:"id" validate:"required,number"`
	Value            string           `json:"value" validate:"required,min=1"`
	AdditionalParams *json.RawMessage `json:"additionalParams"`
}
