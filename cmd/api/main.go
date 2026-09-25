package main

import (
	"go-notification/internal/handler"
	"go-notification/internal/router"
	"go-notification/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	validate := validator.New()

	notificationService := service.NewNotificationService()
	notificationHandler := handler.NewNotificationHandler(notificationService, validate)

	router.Setup(app, notificationHandler)

	app.Listen(":8000")
}
