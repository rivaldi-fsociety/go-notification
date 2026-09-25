package service

import "go-notification/internal/dto"

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) Create(req dto.CreateNotificationRequest) (*dto.CreateNotificationResponse, error) {
	return &dto.CreateNotificationResponse{
		Message: "notification created",
	}, nil
}
