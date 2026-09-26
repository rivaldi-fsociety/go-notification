package helper

import "github.com/gofiber/fiber/v2"

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Success(
	c *fiber.Ctx,
	status int,
	message string,
	data any,
) error {
	return c.Status(status).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithMeta(
	c *fiber.Ctx,
	status int,
	message string,
	data any,
	meta any,
) error {
	return c.Status(status).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(
	c *fiber.Ctx,
	status int,
	message string,
	errors any,
) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
