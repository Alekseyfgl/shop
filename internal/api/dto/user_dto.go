package dto

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,custom_email"`
	Phone    string `json:"phone" validate:"required,min=5"`
	Password string `json:"password" validate:"required,min=3"`
}
