package middlewares

import (
	app_errors "api-gateway/internal/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func (mw *Middlewares) RateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        mw.rateLimitMax,
		Expiration: mw.rateLimitExpiration,
		Storage:    mw.storage,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return app_errors.Respond(c, app_errors.ErrTooManyRequests)
		},
	})
}
