package dto

type CreateSizeRequest struct {
	Title       string  `json:"title" validate:"required,min=1"`
	Description *string `json:"description" validate:"omitempty,min=3,max=1000"`
}

type UpdateSizeRequest struct {
	ID          int     `json:"id" validate:"required,number"`
	Title       string  `json:"title" validate:"required,min=1"`
	Description *string `json:"description" validate:"omitempty,min=3,max=1000"`
}
