package repository

import (
	"context"
	"fmt"
	"go-notification/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type notificationRepository struct {
	collection *mongo.Collection
}

func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &notificationRepository{
		collection: db.Collection("notifications"),
	}
}

func (r *notificationRepository) Create(
	ctx context.Context,
	notification *model.Notification,
) error {
	result, err := r.collection.InsertOne(ctx, notification)
	if err != nil {
		return err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("failed to get inserted notification ID")
	}

	notification.ID = id

	return nil
}
