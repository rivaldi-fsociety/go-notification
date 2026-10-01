package service

import (
	"context"
	"errors"
	"time"

	"go-notification/internal/apperror"
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
	ctx context.Context,
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
		ctx,
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
	ctx context.Context,
	query dto.GetNotificationsQuery,
	userID string,
) ([]dto.NotificationResponse, int64, error) {

	skip := (query.Page - 1) * query.Limit

	notifications, total, err := s.repository.GetAll(
		ctx,
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
	ctx context.Context,
	id string,
	userID string,
) (*dto.NotificationResponse, error) {

	notification, err := s.repository.GetById(
		ctx,
		id,
		userID,
	)

	if err != nil {
		if errors.Is(err, apperror.ErrNotificationNotFound) {
			return nil, apperror.ErrNotificationNotFound
		}

		if errors.Is(err, apperror.ErrInvalidID) {
			return nil, apperror.ErrInvalidID
		}

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
	ctx context.Context,
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
		ctx,
		id,
		notification,
		userID,
	)

	if err != nil {
		if errors.Is(err, apperror.ErrNotificationNotFound) {
			return nil, apperror.ErrNotificationNotFound
		}

		if errors.Is(err, apperror.ErrInvalidID) {
			return nil, apperror.ErrInvalidID
		}

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
	ctx context.Context,
	id string,
	userID string,
) error {
	err := s.repository.Delete(
		ctx,
		id,
		userID,
	)

	if err != nil {
		if errors.Is(err, apperror.ErrNotificationNotFound) {
			return apperror.ErrNotificationNotFound
		}

		if errors.Is(err, apperror.ErrInvalidID) {
			return apperror.ErrInvalidID
		}

		return err
	}

	return nil
}
