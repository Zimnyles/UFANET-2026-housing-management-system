package app_errors

import (
	"errors"
	"net/http"
)

var (
	ErrMissingAuthHeader = errors.New("missing authorization header")
	ErrInvalidAuthFormat = errors.New("invalid authorization format")
	ErrInvalidToken      = errors.New("invalid or expired token")
	ErrInvalidClaims     = errors.New("invalid token claims")

	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")

	ErrClientDial      = errors.New("failed to connect to service")
	ErrTooManyRequests = errors.New("too many requests")
	ErrRequestTimeout  = errors.New("request timeout")
	ErrNotImplemented  = errors.New("not implemented")
	ErrInternal        = errors.New("internal server error")
)

var httpStatus = map[error]int{
	ErrMissingAuthHeader: http.StatusUnauthorized,
	ErrInvalidAuthFormat: http.StatusUnauthorized,
	ErrInvalidToken:      http.StatusUnauthorized,
	ErrInvalidClaims:     http.StatusUnauthorized,
	ErrUnauthorized:      http.StatusUnauthorized,
	ErrForbidden:         http.StatusForbidden,
	ErrNotFound:          http.StatusNotFound,
	ErrBadRequest:        http.StatusBadRequest,
	ErrClientDial:        http.StatusServiceUnavailable,
	ErrTooManyRequests:   http.StatusTooManyRequests,
	ErrRequestTimeout:    http.StatusRequestTimeout,
	ErrNotImplemented:    http.StatusNotImplemented,
	ErrInternal:          http.StatusInternalServerError,
}

func StatusCode(err error) int {
	for target, code := range httpStatus {
		if errors.Is(err, target) {
			return code
		}
	}

	return http.StatusInternalServerError
}
