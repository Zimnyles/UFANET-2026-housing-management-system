package auth_handler

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/dto"
	"net/mail"
)

func validateRegisterRequest(req *dto.RegisterRequest) error {
	if req.Name == "" {
		return app_errors.ErrRegisterNameRequired
	}
	if len(req.Name) < 2 {
		return app_errors.ErrRegisterNameTooShort
	}
	if len(req.Name) > 100 {
		return app_errors.ErrRegisterNameTooLong
	}
	if req.Email == "" {
		return app_errors.ErrRegisterEmailRequired
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return app_errors.ErrRegisterEmailInvalid
	}
	if req.Password == "" {
		return app_errors.ErrRegisterPasswordRequired
	}
	if len(req.Password) < 8 {
		return app_errors.ErrRegisterPasswordTooShort
	}
	if len(req.Password) > 72 {
		return app_errors.ErrRegisterPasswordTooLong
	}
	return nil
}

func validateLoginRequest(req *dto.LoginRequest) error {
	if req.Email == "" {
		return app_errors.ErrLoginEmailRequired
	}
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return app_errors.ErrLoginEmailInvalid
	}
	if req.Password == "" {
		return app_errors.ErrLoginPassRequired
	}
	return nil
}

func validateRefreshRequest(req *dto.RefreshRequest) error {
	if req.RefreshToken == "" {
		return app_errors.ErrRefreshTokenRequired
	}
	return nil
}
