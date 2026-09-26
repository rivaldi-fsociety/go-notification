package router

import (
	"go-notification/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(
	api fiber.Router,
	userHandler *handler.UserHandler,
) {
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
}
