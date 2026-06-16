package middlewares

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(mw.jwtSecret), nil
	})

	if err != nil || claims.ExpiresAt == nil {
		return nil
	}

	ttl := time.Until(claims.ExpiresAt.Time)

	return mw.BlacklistToken(tokenStr, ttl)
}

func (mw *Middlewares) IsBlacklisted(token string) bool {
	val, err := mw.storage.Get(blacklistPrefix + token)

	return err == nil && len(val) > 0
}
