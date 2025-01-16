package dto

type CreateNodeTypeRequest struct {
	Type        string  `json:"type" validate:"required,min=1"`
	Description *string `json:"description" validate:"omitempty,min=3,max=1000"`
}

type UpdateNodeTypeRequest struct {
	ID          int     `json:"id" validate:"required,number"`
	Type        string  `json:"type" validate:"required,min=1"`
	Description *string `json:"description" validate:"omitempty,min=3,max=1000"`
}
