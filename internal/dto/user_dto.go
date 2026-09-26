package dto

import "time"

type RegisterUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Fullname string `json:"full_name" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Fullname  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetUserQuery struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Terms string `query:"terms"`
}
