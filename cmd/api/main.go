package main

import (
	"context"
	"log"

	"go-notification/internal/config"
	"go-notification/internal/database"
	"go-notification/internal/handler"
	"go-notification/internal/helper"
	"go-notification/internal/repository"
	"go-notification/internal/router"
	"go-notification/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// App
	app := fiber.New()
	validate := helper.NewValidator()

	// Database
	mongoClient, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("failed to disconnect MongoDB: %v", err)
		}
	}()

	db := mongoClient.Database(cfg.MongoDatabase)

	// Repositories
	notificationRepository := repository.NewNotificationRepository(db)
	userRepository := repository.NewUserRepository(db)

	// Indexes
	if err := userRepository.CreateIndexes(context.Background()); err != nil {
		log.Fatal(err)
	}

	// Services
	notificationService := service.NewNotificationService(
		notificationRepository,
	)

	userService := service.NewUserService(
		userRepository,
		cfg.JWTSecret,
		cfg.JWTExpiresIn,
	)

	// Handlers
	notificationHandler := handler.NewNotificationHandler(
		notificationService,
		validate,
	)

	userHandler := handler.NewUserHandler(
		userService,
		validate,
	)

	// Router
	router.Setup(
		app,
		notificationHandler,
		userHandler,
		cfg.JWTSecret,
		cfg.RequestTimeout,
	)

	// Server
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
