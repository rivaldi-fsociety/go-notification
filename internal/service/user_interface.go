package service

import "go-notification/internal/dto"

type UserService interface {
	Register(req dto.RegisterUserRequest) (*dto.UserResponse, error)
	Login(req dto.LoginUserRequest) (*dto.LoginResponse, error)
}
