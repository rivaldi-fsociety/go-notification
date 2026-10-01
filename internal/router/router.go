package router

import (
	"go-notification/internal/handler"
	"go-notification/internal/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Setup(
	app *fiber.App,
	notificationHandler *handler.NotificationHandler,
	userHandler *handler.UserHandler,
	jwtSecret string,
	requestTimeout time.Duration,
) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public
	v1.Get("/health", handler.HealthCheck)

	SetupAuthRoutes(v1, userHandler)

	// Protected
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	timeoutMiddleware := middleware.TimeoutMiddleware(requestTimeout)

	SetupNotificationRouter(
		v1,
		notificationHandler,
		authMiddleware,
		timeoutMiddleware,
	)
}
