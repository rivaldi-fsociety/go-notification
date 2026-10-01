package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"go-notification/internal/dto"
	"go-notification/internal/helper"
	"go-notification/internal/service"
)

type UserHandler struct {
	service   service.UserService
	validator *validator.Validate
}

func NewUserHandler(
	service service.UserService,
	validator *validator.Validate,
) *UserHandler {
	return &UserHandler{
		service:   service,
		validator: validator,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var req dto.RegisterUserRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Error(
			c,
			http.StatusBadRequest,
			"Invalid request body",
			nil,
		)
	}

	if err := h.validator.Struct(req); err != nil {
		return helper.Error(
			c,
			http.StatusBadRequest,
			"Validation failed",
			helper.ValidationErrors(err),
		)
	}

	user, err := h.service.Register(ctx, req)
	if err != nil {
		return helper.Error(
			c,
			http.StatusInternalServerError,
			err.Error(),
			nil,
		)
	}

	return helper.Success(
		c,
		http.StatusCreated,
		"User registered successfully",
		user,
	)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()
	var req dto.LoginUserRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(
			c,
			http.StatusBadRequest,
			"Invalid request body",
			nil,
		)
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		return helper.Error(
			c,
			http.StatusBadRequest,
			"Validation failed",
			helper.ValidationErrors(err),
		)
	}

	// Login
	result, err := h.service.Login(ctx, req)
	if err != nil {
		return handleError(c, err)
	}

	return helper.Success(
		c,
		http.StatusOK,
		"Login successful",
		result,
	)
}
