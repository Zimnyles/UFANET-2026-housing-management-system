package middlewares

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/constants"
	"github.com/gofiber/fiber/v2"
)

func (mw *Middlewares) RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(c *fiber.Ctx) error {
		role, ok := c.Locals(constants.LocalRole).(string)
		if !ok || role == "" {
			return app_errors.Respond(c, app_errors.ErrUnauthorized)
		}

		if _, ok := allowed[role]; !ok {
			return app_errors.Respond(c, app_errors.ErrForbidden)
		}

		return c.Next()
	}
}
