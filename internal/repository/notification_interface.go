package repository

import (
	"context"
	"go-notification/internal/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
}
