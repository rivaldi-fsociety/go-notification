package service

import "go-notification/internal/dto"

type NotificationService interface {
	Create(req dto.CreateNotificationRequest) (*dto.CreateNotificationResponse, error)
}
