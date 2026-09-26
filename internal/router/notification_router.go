package router

import (
	"go-notification/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupNotificationRouter(
	api fiber.Router,
	notificationHandler *handler.NotificationHandler,
) {
	notifications := api.Group("/notifications")
	notifications.Post("/", notificationHandler.Create)
	notifications.Get("/", notificationHandler.GetAll)
	notifications.Get("/:id", notificationHandler.Get)
	notifications.Put("/:id", notificationHandler.Update)
	notifications.Delete("/:id", notificationHandler.Delete)
}
