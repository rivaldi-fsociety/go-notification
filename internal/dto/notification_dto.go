package dto

type CreateNotificationRequest struct {
	UserID  string `json:"user_id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type CreateNotificationResponse struct {
	Message string `json:"message"`
}
