package service

import (
	"context"
	"time"

	"go-notification/internal/dto"
	"go-notification/internal/model"
	"go-notification/internal/repository"
)

type notificationService struct {
	repository repository.NotificationRepository
}

func NewNotificationService(
	repository repository.NotificationRepository,
) NotificationService {
	return &notificationService{
		repository: repository,
	}
}

func (s *notificationService) Create(
	req dto.CreateNotificationRequest,
	userID string,
) (*dto.NotificationResponse, error) {

	notification := &model.Notification{
		UserID:    userID,
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

	return &dto.NotificationResponse{
		ID:        notification.ID.Hex(),
		UserID:    notification.UserID,
		Title:     notification.Title,
		Message:   notification.Message,
		CreatedAt: notification.CreatedAt,
	}, nil
}

func (s *notificationService) GetAll(
	query dto.GetNotificationsQuery,
	userID string,
) ([]dto.NotificationResponse, int64, error) {

	skip := (query.Page - 1) * query.Limit

	notifications, total, err := s.repository.GetAll(
		context.Background(),
		repository.GetNotificationsParams{
			Limit:  query.Limit,
			Skip:   skip,
			Terms:  query.Terms,
			UserID: userID,
		},
	)

	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.NotificationResponse, 0, len(notifications))

	for _, notification := range notifications {
		result = append(result, dto.NotificationResponse{
			ID:        notification.ID.Hex(),
			UserID:    notification.UserID,
			Title:     notification.Title,
			Message:   notification.Message,
			CreatedAt: notification.CreatedAt,
		})
	}

	return result, total, nil
}

func (s *notificationService) Get(
	id string,
	userID string,
) (*dto.NotificationResponse, error) {

	notification, err := s.repository.GetById(
		context.Background(),
		id,
		userID,
	)

	if err != nil {
		return nil, err
	}

	return &dto.NotificationResponse{
		ID:        notification.ID.Hex(),
		UserID:    notification.UserID,
		Title:     notification.Title,
		Message:   notification.Message,
		CreatedAt: notification.CreatedAt,
	}, nil
}

func (s *notificationService) Update(
	req dto.CreateNotificationRequest,
	id string,
	userID string,
) (*dto.NotificationResponse, error) {

	notification := &model.Notification{
		UserID:    userID,
		Title:     req.Title,
		Message:   req.Message,
		UpdatedAt: time.Now(),
	}

	updatedNotification, err := s.repository.Update(
		context.Background(),
		id,
		notification,
		userID,
	)

	if err != nil {
		return nil, err
	}

	return &dto.NotificationResponse{
		ID:        updatedNotification.ID.Hex(),
		UserID:    updatedNotification.UserID,
		Title:     updatedNotification.Title,
		Message:   updatedNotification.Message,
		CreatedAt: updatedNotification.CreatedAt,
		UpdatedAt: updatedNotification.UpdatedAt,
	}, nil
}

func (s *notificationService) Delete(
	id string,
	userID string,
) error {
	err := s.repository.Delete(
		context.Background(),
		id,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}
