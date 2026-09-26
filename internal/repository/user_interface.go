package repository

import (
	"context"
	"go-notification/internal/model"
)

type GetUserParams struct {
	Limit int
	Skip  int
	Terms string
}

type UserRepository interface {
	CreateIndexes(ctx context.Context) error

	Create(ctx context.Context, user *model.User) error
	GetAll(ctx context.Context, params GetUserParams) ([]model.User, int64, error)
	GetById(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id string, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id string) error
}
