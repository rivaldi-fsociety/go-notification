package service

import (
	"context"
	"errors"
	"fmt"
	"go-notification/internal/dto"
	"go-notification/internal/model"
	"go-notification/internal/repository"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repository repository.UserRepository
}

func NewUserService(
	repository repository.UserRepository,
) UserService {
	return &userService{
		repository: repository,
	}
}

func (s *userService) Register(
	req dto.RegisterUserRequest,
) (*dto.UserResponse, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))

	existingUser, err := s.repository.FindByEmail(
		context.Background(),
		email,
	)

	if err != nil {
		if !errors.Is(err, repository.ErrUserNotFound) {
			return nil, err
		}
	}

	if existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	user := &model.User{
		Email:     email,
		Password:  string(hashedPassword),
		Fullname:  req.Fullname,
		Phone:     req.Phone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Create(
		context.Background(),
		user,
	); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Fullname:  user.Fullname,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
