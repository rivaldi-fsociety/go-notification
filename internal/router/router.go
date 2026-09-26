package router

import (
	"go-notification/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func Setup(
	app *fiber.App,
	notificationHandler *handler.NotificationHandler,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/health", handler.HealthCheck)
	notifications := v1.Group("/notifications")
	notifications.Post("/", notificationHandler.Create)
	notifications.Get("/", notificationHandler.GetAll)
}
