package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	MongoURI      string
	MongoDatabase string
	JWTSecret     string
	JWTExpiresIn  time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiresIn, err := time.ParseDuration(
		getEnv("JWT_EXPIRES_IN", "24h"),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRES_IN: %w", err)
	}

	config := &Config{
		AppPort:       getEnv("APP_PORT", "8000"),
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "go_notification"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpiresIn:  jwtExpiresIn,
	}

	if config.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return config, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
