package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Notification struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string        `bson:"user_id" json:"user_id"`
	Title     string        `bson:"title" json:"title"`
	Message   string        `bson:"message" json:"message"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
