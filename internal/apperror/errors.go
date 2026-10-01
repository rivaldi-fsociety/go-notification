package apperror

import "errors"

var (
	ErrInvalidID            = errors.New("invalid id")
	ErrNotificationNotFound = errors.New("notification not found")

	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
