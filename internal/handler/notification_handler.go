package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"go-notification/internal/dto"
	"go-notification/internal/helper"
	"go-notification/internal/service"
)

type NotificationHandler struct {
	service   service.NotificationService
	validator *validator.Validate
}

func NewNotificationHandler(
	service service.NotificationService,
	validator *validator.Validate,
) *NotificationHandler {
	return &NotificationHandler{
		service:   service,
		validator: validator,
	}
}

func (h *NotificationHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateNotificationRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Error(
			c,
			fiber.StatusBadRequest,
			"invalid request body",
			nil,
		)
	}

	if err := h.validator.Struct(req); err != nil {
		return helper.Error(
			c,
			fiber.StatusBadRequest,
			"validation failed",
			helper.ValidationErrors(err),
		)
	}

	result, err := h.service.Create(req)
	if err != nil {
		return helper.Error(
			c,
			fiber.StatusInternalServerError,
			"failed to create notification",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusCreated,
		"notification created",
		result,
	)
}
