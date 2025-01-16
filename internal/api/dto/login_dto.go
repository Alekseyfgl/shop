package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,custom_email"`
	Password string `json:"password" validate:"required,min=3"`
}
