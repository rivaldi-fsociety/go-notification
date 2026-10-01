package repository

import (
	"context"
	"errors"
	"fmt"
	"go-notification/internal/apperror"
	"go-notification/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (r *notificationRepository) GetAll(
	ctx context.Context,
	params GetNotificationsParams,
) ([]model.Notification, int64, error) {

	filter := bson.M{
		"user_id": params.UserID,
	}

	if params.Terms != "" {
		filter = bson.M{
			"$or": bson.A{
				bson.M{
					"title": bson.M{
						"$regex":   params.Terms,
						"$options": "i",
					},
				},
				bson.M{
					"message": bson.M{
						"$regex":   params.Terms,
						"$options": "i",
					},
				},
			},
		}
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	options := options.Find().
		SetSkip(int64(params.Skip)).
		SetLimit(int64(params.Limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var notifications []model.Notification

	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *notificationRepository) GetById(
	ctx context.Context,
	id string,
	userID string,
) (*model.Notification, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	var notification model.Notification

	err = r.collection.FindOne(ctx, filter).Decode(&notification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrNotificationNotFound
		}
		return nil, err
	}

	return &notification, nil
}

func (r *notificationRepository) Update(
	ctx context.Context,
	id string,
	notification *model.Notification,
	userID string,
) (*model.Notification, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	update := bson.M{
		"$set": bson.M{
			"title":      notification.Title,
			"message":    notification.Message,
			"updated_at": notification.UpdatedAt,
		},
	}

	options := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var updatedNotification model.Notification

	err = r.collection.
		FindOneAndUpdate(ctx, filter, update, options).
		Decode(&updatedNotification)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrNotificationNotFound
		}

		return nil, err
	}

	return &updatedNotification, nil
}

func (r *notificationRepository) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrInvalidID
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{
		"_id":     objectID,
		"user_id": userID,
	})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return apperror.ErrNotificationNotFound
	}

	return nil
}
