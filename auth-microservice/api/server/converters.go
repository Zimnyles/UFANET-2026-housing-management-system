package server

import (
	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"

	"auth-service/infra/models/domain"
)

func ConvertFromProtoToRegisterRequest(req *authpb.RegisterRequest) *domain.RegisterRequest {
	return &domain.RegisterRequest{
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		AdminCode: req.GetAdminCode(),
	}
}

func ConvertFromProtoToLoginRequest(req *authpb.LoginRequest) *domain.LoginRequest {
	return &domain.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func ConvertFromProtoToRefreshRequest(req *authpb.RefreshRequest) *domain.RefreshRequest {
	return &domain.RefreshRequest{
		RefreshToken: req.GetRefreshToken(),
	}
}

func ConvertFromProtoToLogoutRequest(req *authpb.LogoutRequest) *domain.LogoutRequest {
	return &domain.LogoutRequest{
		RefreshToken: req.GetRefreshToken(),
	}
}

func ConvertFromDomainToAuthResponse(d *domain.AuthResult) *authpb.AuthResponse {
	return &authpb.AuthResponse{
		UserId:       d.UserID,
		AccessToken:  d.AccessToken,
		RefreshToken: d.RefreshToken,
	}
}

func ConvertFromDomainToRefreshResponse(d *domain.RefreshResult) *authpb.RefreshResponse {
	return &authpb.RefreshResponse{
		AccessToken: d.AccessToken,
	}
}
