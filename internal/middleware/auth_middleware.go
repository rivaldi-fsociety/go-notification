package middleware

import (
	"net/http"
	"strings"

	"go-notification/internal/helper"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return helper.Error(
				c,
				http.StatusUnauthorized,
				"Authorization header is required",
				nil,
			)
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			return helper.Error(
				c,
				http.StatusUnauthorized,
				"Invalid authorization header",
				nil,
			)
		}

		tokenString := parts[1]

		claims, err := helper.ParseJWT(
			tokenString,
			jwtSecret,
		)
		if err != nil {
			return helper.Error(
				c,
				http.StatusUnauthorized,
				"Invalid or expired token",
				nil,
			)
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}
