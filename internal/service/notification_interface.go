package service

import "go-notification/internal/dto"

type NotificationService interface {
	Create(req dto.CreateNotificationRequest) (*dto.NotificationResponse, error)
	GetAll(query dto.GetNotificationsQuery) ([]dto.NotificationResponse, int64, error)
	Get(id string) (*dto.NotificationResponse, error)
	Update(req dto.CreateNotificationRequest, id string) (*dto.NotificationResponse, error)
}
