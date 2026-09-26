package main

import (
	"context"
	"log"

	"go-notification/internal/database"
	"go-notification/internal/handler"
	"go-notification/internal/helper"
	"go-notification/internal/repository"
	"go-notification/internal/router"
	"go-notification/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	validate := helper.NewValidator()

	// Database
	mongoClient, err := database.ConnectMongoDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("failed to disconnect MongoDB: %v", err)
		}
	}()

	db := mongoClient.Database("go_notification")

	// Dependencies
	notificationRepository := repository.NewNotificationRepository(db)
	notificationService := service.NewNotificationService(notificationRepository)
	notificationHandler := handler.NewNotificationHandler(
		notificationService,
		validate,
	)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService, validate)

	if err := userRepository.CreateIndexes(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Router
	router.Setup(app, notificationHandler, userHandler)

	// Server
	if err := app.Listen(":8000"); err != nil {
		log.Fatal(err)
	}
}
