package handler

import (
	"errors"
	"log"
	"net/http"

	"go-notification/internal/apperror"
	"go-notification/internal/helper"

	"github.com/gofiber/fiber/v2"
)

func handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, apperror.ErrInvalidID):
		return helper.Error(
			c,
			http.StatusBadRequest,
			"Invalid ID",
			nil,
		)

	case errors.Is(err, apperror.ErrNotificationNotFound):
		return helper.Error(
			c,
			http.StatusNotFound,
			"Notification not found",
			nil,
		)

	case errors.Is(err, apperror.ErrEmailAlreadyExists):
		return helper.Error(
			c,
			http.StatusConflict,
			"Email already registered",
			nil,
		)

	case errors.Is(err, apperror.ErrInvalidCredentials):
		return helper.Error(
			c,
			http.StatusUnauthorized,
			"Invalid email or password",
			nil,
		)

	default:
		log.Printf(
			"request_id=%v unexpected_error=%v",
			c.Locals("requestid"),
			err,
		)

		return helper.Error(
			c,
			http.StatusInternalServerError,
			"Internal server error",
			nil,
		)
	}
}
