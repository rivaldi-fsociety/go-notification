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

func (h *NotificationHandler) GetAll(c *fiber.Ctx) error {
	query := dto.GetNotificationsQuery{
		Page:  1,
		Limit: 10,
	}

	if err := c.QueryParser(&query); err != nil {
		return helper.Error(
			c,
			fiber.StatusBadRequest,
			"invalid query parameters",
			nil,
		)
	}

	notifications, total, err := h.service.GetAll(query)
	if err != nil {
		return helper.Error(
			c,
			fiber.StatusInternalServerError,
			"failed to get notifications",
			nil,
		)
	}

	pagination := helper.Pagination{
		Page:       query.Page,
		Limit:      query.Limit,
		Total:      total,
		TotalPages: int((total + int64(query.Limit) - 1) / int64(query.Limit)),
	}

	return helper.SuccessWithMeta(
		c,
		fiber.StatusOK,
		"notifications retrieved",
		notifications,
		pagination,
	)
}

func (h *NotificationHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")

	result, err := h.service.Get(id)
	if err != nil {
		return helper.Error(
			c,
			fiber.StatusInternalServerError,
			"failed to get notification",
			nil,
		)
	}

	if result == nil {
		return helper.Error(
			c,
			fiber.StatusNotFound,
			"notification not found",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"notification retrieved",
		result,
	)
}

func (h *NotificationHandler) Update(c *fiber.Ctx) error {
	var req dto.CreateNotificationRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Error(
			c,
			fiber.StatusBadRequest,
			"invalid request body",
			nil,
		)
	}

	id := c.Params("id")

	result, err := h.service.Update(req, id)
	if err != nil {
		return helper.Error(
			c,
			fiber.StatusInternalServerError,
			"failed to update notification",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"notification updated",
		result,
	)
}
