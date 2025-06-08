package dto

type OrderDTO struct {
	NodeId int     `json:"nodeId" validate:"required,number"`
	Size   *string `json:"description" validate:"omitempty,min=3,max=1000"`
	Amount int     `json:"amount" validate:"required,number"`
}
