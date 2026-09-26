package service

import "go-notification/internal/dto"

type NotificationService interface {
	Create(req dto.CreateNotificationRequest) (*dto.NotificationResponse, error)
	GetAll(query dto.GetNotificationsQuery) ([]dto.NotificationResponse, int64, error)
}
