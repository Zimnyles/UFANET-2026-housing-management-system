package server

import (
	infra_errors "auth-service/infra/errors"
	"net/mail"

	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"
)

func ValidateRegisterRequest(req *authpb.RegisterRequest) error {
	if req.GetName() == "" {
		return infra_errors.ErrNameRequired
	}
	if len(req.GetName()) < 2 {
		return infra_errors.ErrNameTooShort
	}
	if len(req.GetName()) > 100 {
		return infra_errors.ErrNameTooLong
	}
	if req.GetEmail() == "" {
		return infra_errors.ErrEmailRequired
	}
	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return infra_errors.ErrEmailInvalid
	}
	if req.GetPassword() == "" {
		return infra_errors.ErrPasswordRequired
	}
	if len(req.GetPassword()) < 8 {
		return infra_errors.ErrPasswordTooShort
	}
	if len(req.GetPassword()) > 72 {
		return infra_errors.ErrPasswordTooLong
	}
	if len(req.GetAdminCode()) > 72 {
		return infra_errors.ErrInvalidAdminCode
	}
	return nil
}

func ValidateLoginRequest(req *authpb.LoginRequest) error {
	if req.GetEmail() == "" {
		return infra_errors.ErrEmailRequired
	}
	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return infra_errors.ErrEmailInvalid
	}
	if req.GetPassword() == "" {
		return infra_errors.ErrPasswordRequired
	}
	return nil
}

func ValidateRefreshRequest(req *authpb.RefreshRequest) error {
	if req.GetRefreshToken() == "" {
		return infra_errors.ErrRefreshRequired
	}
	return nil
}

func ValidateLogoutRequest(req *authpb.LogoutRequest) error {
	if req.GetRefreshToken() == "" {
		return infra_errors.ErrRefreshRequired
	}
	return nil
}
