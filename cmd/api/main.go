package main

import (
	"go-notification/internal/handler"
	"go-notification/internal/helper"
	"go-notification/internal/router"
	"go-notification/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	validate := helper.NewValidator()

	notificationService := service.NewNotificationService()
	notificationHandler := handler.NewNotificationHandler(notificationService, validate)

	router.Setup(app, notificationHandler)

	app.Listen(":8000")
}
