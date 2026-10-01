package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TimeoutMiddleware(duration time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(
			c.UserContext(),
			duration,
		)
		defer cancel()

		c.SetUserContext(ctx)

		return c.Next()
	}
}
