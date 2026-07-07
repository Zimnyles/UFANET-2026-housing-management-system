package errors

import "errors"

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrInvalidToken            = errors.New("invalid token")
	ErrTokenExpired            = errors.New("token expired")
	ErrTokenNotFound           = errors.New("refresh token not found")
	ErrInternal                = errors.New("internal server error")
	ErrUnexpectedSigningMethod = errors.New("unexpected JWT signing method")
)

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailInvalid     = errors.New("invalid email format")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password min length 8")
	ErrPasswordTooLong  = errors.New("password max length 72")
	ErrRefreshRequired  = errors.New("refresh_token is required")
	ErrInvalidAdminCode = errors.New("invalid admin code")
)
