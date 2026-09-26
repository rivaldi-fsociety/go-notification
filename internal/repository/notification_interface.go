package repository

import (
	"context"
	"go-notification/internal/model"
)

type GetNotificationsParams struct {
	Limit int
	Skip  int
	Terms string
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	GetAll(ctx context.Context, params GetNotificationsParams) ([]model.Notification, int64, error)
	GetById(ctx context.Context, id string) (*model.Notification, error)
	Update(ctx context.Context, id string, notification *model.Notification) (*model.Notification, error)
	Delete(ctx context.Context, id string) error
}
