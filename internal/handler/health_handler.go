package handler

import (
	"go-notification/internal/helper"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(c *fiber.Ctx) error {
	return helper.Success(
		c,
		fiber.StatusOK,
		"go-notification running",
		nil,
	)
}
