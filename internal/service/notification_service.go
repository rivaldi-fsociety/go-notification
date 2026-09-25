package service

import (
	"context"
	"time"

	"go-notification/internal/dto"
	"go-notification/internal/model"
	"go-notification/internal/repository"
)

type NotificationService struct {
	repository *repository.NotificationRepository
}

func NewNotificationService(
	repository *repository.NotificationRepository,
) *NotificationService {
	return &NotificationService{
		repository: repository,
	}
}

func (s *NotificationService) Create(
	req dto.CreateNotificationRequest,
) (*dto.CreateNotificationResponse, error) {

	notification := &model.Notification{
		UserID:    req.UserID,
		Title:     req.Title,
		Message:   req.Message,
		CreatedAt: time.Now(),
	}

	err := s.repository.Create(
		context.Background(),
		notification,
	)

	if err != nil {
		return nil, err
	}

	return &dto.CreateNotificationResponse{
		ID:        notification.ID.Hex(),
		UserID:    notification.UserID,
		Title:     notification.Title,
		Message:   notification.Message,
		CreatedAt: notification.CreatedAt,
	}, nil
}
