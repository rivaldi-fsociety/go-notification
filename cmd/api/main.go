package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-notification/internal/cache"
	"go-notification/internal/config"
	"go-notification/internal/database"
	"go-notification/internal/handler"
	"go-notification/internal/helper"
	"go-notification/internal/middleware"
	"go-notification/internal/repository"
	"go-notification/internal/router"
	"go-notification/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// App
	app := fiber.New()
	app.Use(requestid.New())
	app.Use(middleware.LoggerMiddleware())
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

	// Redis
	redisClient, err := database.ConnectRedis(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("failed to close redis: %v", err)
		}
	}()

	log.Println("redis connected successfully")

	notificationCache := cache.NewRedisCache(redisClient)

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
		notificationCache,
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

	shutdownCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Server
	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	<-shutdownCtx.Done()

	log.Println("shutdown signal received")

	shutdownTimeoutCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownTimeoutCtx); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	} else {
		log.Println("server shutdown complete")
	}
}
