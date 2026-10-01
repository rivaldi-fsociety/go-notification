package service

import "go-notification/internal/dto"

type NotificationService interface {
	Create(req dto.CreateNotificationRequest, userID string) (*dto.NotificationResponse, error)
	GetAll(query dto.GetNotificationsQuery, userID string) ([]dto.NotificationResponse, int64, error)
	Get(id string, userID string) (*dto.NotificationResponse, error)
	Update(req dto.CreateNotificationRequest, id string, userID string) (*dto.NotificationResponse, error)
	Delete(id string, userID string) error
}
