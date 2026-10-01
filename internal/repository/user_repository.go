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

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) CreateIndexes(
	ctx context.Context,
) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().
			SetUnique(true),
	}

	_, err := r.collection.Indexes().CreateOne(
		ctx,
		indexModel,
	)

	return err
}

func (r *userRepository) Create(
	ctx context.Context,
	user *model.User,
) error {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("failed to get inserted user ID")
	}

	user.ID = id

	return nil
}

func (r *userRepository) GetAll(
	ctx context.Context,
	params GetUserParams,
) ([]model.User, int64, error) {

	filter := bson.M{}

	if params.Terms != "" {
		filter = bson.M{
			"$or": bson.A{
				bson.M{
					"email": bson.M{
						"$regex":   params.Terms,
						"$options": "i",
					},
				},
				bson.M{
					"full_name": bson.M{
						"$regex":   params.Terms,
						"$options": "i",
					},
				},
				bson.M{
					"phone": bson.M{
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

	var users []model.User

	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetById(
	ctx context.Context,
	id string,
) (*model.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	filter := bson.M{"_id": objectID}

	var user model.User

	err = r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {

	filter := bson.M{
		"email": email,
	}

	var user model.User

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(
	ctx context.Context,
	id string,
	user *model.User,
) (*model.User, error) {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	filter := bson.M{
		"_id": objectID,
	}

	update := bson.M{
		"$set": bson.M{
			"full_name":  user.Fullname,
			"phone":      user.Phone,
			"updated_at": user.UpdatedAt,
		},
	}

	options := options.FindOneAndUpdate().
		SetReturnDocument(options.After)

	var updatedUser model.User

	err = r.collection.
		FindOneAndUpdate(ctx, filter, update, options).
		Decode(&updatedUser)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("user not found")
		}

		return nil, err
	}

	return &updatedUser, nil
}

func (r *userRepository) Delete(
	ctx context.Context,
	id string,
) error {

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid user id")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{
		"_id": objectID,
	})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
