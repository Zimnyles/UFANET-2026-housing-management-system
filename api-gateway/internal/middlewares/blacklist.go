package middlewares

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	app_errors "api-gateway/internal/errors"
)

const blacklistPrefix = "bl:"

func (mw *Middlewares) BlacklistToken(token string, ttl time.Duration) error {
	if ttl <= 0 {
		return nil
	}

	return mw.storage.Set(blacklistPrefix+token, []byte("1"), ttl)
}

func (mw *Middlewares) BlacklistRawJWT(tokenStr string) error {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", app_errors.ErrUnexpectedSigningMethod, t.Header["alg"])
		}

		return []byte(mw.jwtSecret), nil
	})
	if err != nil {
		return fmt.Errorf("%w: %w", app_errors.ErrInvalidToken, err)
	}

	if claims.ExpiresAt == nil {
		return app_errors.ErrInvalidClaims
	}

	ttl := time.Until(claims.ExpiresAt.Time)

	return mw.BlacklistToken(tokenStr, ttl)
}

func (mw *Middlewares) IsBlacklisted(token string) bool {
	val, err := mw.storage.Get(blacklistPrefix + token)

	return err == nil && len(val) > 0
}
