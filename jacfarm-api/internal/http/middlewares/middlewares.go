package middlewares

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func ServiceAuthMiddleware(apiKey string) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		token := c.Get("Authorization")

		if token != apiKey {
			return c.SendStatus(http.StatusUnauthorized)
		}

		return c.Next()
	}
}
