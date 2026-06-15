package middlewares

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/constants"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (mw *Middlewares) JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return app_errors.Respond(c, app_errors.ErrMissingAuthHeader)
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return app_errors.Respond(c, app_errors.ErrInvalidAuthFormat)
		}

		token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(mw.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return app_errors.Respond(c, app_errors.ErrInvalidToken)
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			return app_errors.Respond(c, app_errors.ErrInvalidClaims)
		}

		c.Locals(constants.LocalUserID, claims.UserID)
		c.Locals(constants.LocalRole, claims.Role)
		return c.Next()
	}
}
