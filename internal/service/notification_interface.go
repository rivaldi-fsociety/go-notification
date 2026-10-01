package service

import (
	"context"
	"go-notification/internal/dto"
)

type NotificationService interface {
	Create(ctx context.Context, req dto.CreateNotificationRequest, userID string) (*dto.NotificationResponse, error)
	GetAll(ctx context.Context, query dto.GetNotificationsQuery, userID string) ([]dto.NotificationResponse, int64, error)
	Get(ctx context.Context, id string, userID string) (*dto.NotificationResponse, error)
	Update(ctx context.Context, req dto.CreateNotificationRequest, id string, userID string) (*dto.NotificationResponse, error)
	Delete(ctx context.Context, id string, userID string) error
}
