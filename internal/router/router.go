package router

import (
	"go-notification/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func Setup(
	app *fiber.App,
	notificationHandler *handler.NotificationHandler,
	userHandler *handler.UserHandler,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", handler.HealthCheck)

	SetupNotificationRouter(v1, notificationHandler)
	SetupAuthRoutes(v1, userHandler)
}
