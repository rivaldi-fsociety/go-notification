package dto

import "time"

type CreateNotificationRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type NotificationResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type GetNotificationsQuery struct {
	Page  int    `query:"page"`
	Limit int    `query:"limit"`
	Terms string `query:"terms"`
}
