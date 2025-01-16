package dto

type CreateNodeRequest struct {
	Title       string  `json:"title" validate:"required,min=1"`
	NodeTypeId  int     `json:"nodeTypeId" validate:"required,number"`
	Description *string `json:"description" validate:"omitempty,min=3,max=1000"`
}

type UpdateNodeRequest struct {
	ID          int     `json:"id" validate:"required,number"`
	Title       string  `json:"title" validate:"required,min=1"`
	NodeTypeId  int     `json:"nodeTypeId" validate:"required,number"`
	Description *string `json:"description" validate:"omitempty,min=3,max=1000"`
}
