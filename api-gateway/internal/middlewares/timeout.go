package middlewares

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

func (mw *Middlewares) Timeout() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), mw.requestTimeout)
		defer cancel()
		c.SetUserContext(ctx)

		return c.Next()
	}
}
