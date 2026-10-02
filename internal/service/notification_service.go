package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go-notification/internal/cache"
	"go-notification/internal/dto"
	"go-notification/internal/model"
	"go-notification/internal/repository"
)

type notificationService struct {
	repository repository.NotificationRepository
	cache      cache.Cache
}

func NewNotificationService(
	repository repository.NotificationRepository,
	cache cache.Cache,
) NotificationService {
	return &notificationService{
		repository: repository,
		cache:      cache,
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

	cacheKey := fmt.Sprintf("notification:%s:%s", userID, id)

	cached, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		log.Printf("failed to get notification from cache: %v", err)
	} else if cached != "" {
		var response dto.NotificationResponse

		if err := json.Unmarshal([]byte(cached), &response); err != nil {
			log.Printf("failed to unmarshal cached notification: %v", err)
		} else {
			return &response, nil
		}
	}

	notification, err := s.repository.GetById(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	response := dto.NotificationResponse{
		ID:        notification.ID.Hex(),
		UserID:    notification.UserID,
		Title:     notification.Title,
		Message:   notification.Message,
		CreatedAt: notification.CreatedAt,
	}

	data, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(
		ctx,
		cacheKey,
		string(data),
		5*time.Minute,
	); err != nil {
		log.Printf("failed to cache notification: %v", err)
	}

	return &response, nil
}

func (s *notificationService) Update(
	ctx context.Context,
	req dto.CreateNotificationRequest,
	id string,
	userID string,
) (*dto.NotificationResponse, error) {

	notification := &model.Notification{
		Title:   req.Title,
		Message: req.Message,
	}

	updatedNotification, err := s.repository.Update(
		ctx,
		id,
		notification,
		userID,
	)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("notification:%s:%s", userID, id)

	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		log.Printf("failed to invalidate notification cache: %v", err)
	}

	response := dto.NotificationResponse{
		ID:        updatedNotification.ID.Hex(),
		UserID:    updatedNotification.UserID,
		Title:     updatedNotification.Title,
		Message:   updatedNotification.Message,
		CreatedAt: updatedNotification.CreatedAt,
	}

	return &response, nil
}

func (s *notificationService) Delete(
	ctx context.Context,
	userID string,
	id string,
) error {

	err := s.repository.Delete(ctx, id, userID)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("notification:%s:%s", userID, id)

	if err := s.cache.Delete(ctx, cacheKey); err != nil {
		log.Printf("failed to invalidate notification cache: %v", err)
	}

	return nil
}
