package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (mw *Middlewares) requestLogger(logger *zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		requestLogger := logger.With().
			Str("request_id", requestID).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Logger()

		c.Set("X-Request-ID", requestID)

		err := c.Next()

		requestLogger.Info().
			Int("status", c.Response().StatusCode()).
			Dur("duration", time.Since(start)).
			Msg("request completed")

		return err
	}
}
