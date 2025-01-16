package dto

type CreateCharDefValueRequest struct {
	CharacteristicId int    `json:"characteristicId" validate:"required,number"`
	Value            string `json:"value" validate:"required,min=1"`
}

type UpdateCharDefValueRequest struct {
	ID    int    `json:"id" validate:"required,number"`
	Value string `json:"value" validate:"required,min=1"`
}
