package service

import (
	"context"
	"go-notification/internal/dto"
)

type UserService interface {
	Register(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginUserRequest) (*dto.LoginResponse, error)
}
