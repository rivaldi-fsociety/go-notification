package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		requestID := c.Locals("requestid")

		log.Printf(
			"request_id=%v method=%s path=%s status=%d duration=%s",
			requestID,
			c.Method(),
			c.Path(),
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}
