package router

import (
	"go-notification/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupNotificationRouter(
	router fiber.Router,
	notificationHandler *handler.NotificationHandler,
	authMiddleware fiber.Handler,
	timeoutMiddleware fiber.Handler,
) {
	notifications := router.Group(
		"/notifications",
		authMiddleware,
		timeoutMiddleware,
	)

	notifications.Post("/", notificationHandler.Create)
	notifications.Get("/", notificationHandler.GetAll)
	notifications.Get("/:id", notificationHandler.Get)
	notifications.Put("/:id", notificationHandler.Update)
	notifications.Delete("/:id", notificationHandler.Delete)
}
