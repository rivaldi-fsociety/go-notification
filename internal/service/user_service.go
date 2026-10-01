package service

import (
	"context"
	"errors"
	"go-notification/internal/apperror"
	"go-notification/internal/dto"
	"go-notification/internal/helper"
	"go-notification/internal/model"
	"go-notification/internal/repository"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repository   repository.UserRepository
	jwtSecret    string
	jwtExpiresIn time.Duration
}

func NewUserService(
	repository repository.UserRepository,
	jwtSecret string,
	jwtExpiresIn time.Duration,
) UserService {
	return &userService{
		repository:   repository,
		jwtSecret:    jwtSecret,
		jwtExpiresIn: jwtExpiresIn,
	}
}

func (s *userService) Register(
	ctx context.Context,
	req dto.RegisterUserRequest,
) (*dto.UserResponse, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))

	existingUser, err := s.repository.FindByEmail(
		ctx,
		email,
	)

	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, apperror.ErrEmailAlreadyExists
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

func (s *userService) Login(
	ctx context.Context,
	req dto.LoginUserRequest,
) (*dto.LoginResponse, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.repository.FindByEmail(
		ctx,
		email,
	)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			return nil, apperror.ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, apperror.ErrInvalidCredentials
	}

	token, err := helper.GenerateJWT(
		user.ID.Hex(),
		user.Email,
		s.jwtSecret,
		s.jwtExpiresIn,
	)
	if err != nil {
		return nil, err
	}

	userResponse := dto.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Fullname:  user.Fullname,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &dto.LoginResponse{
		Token: token,
		User:  userResponse,
	}, nil
}
