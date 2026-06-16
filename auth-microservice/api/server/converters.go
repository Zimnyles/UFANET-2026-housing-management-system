package server

import (
	"auth-service/infra/models/dto"

	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"
)

func ConvertFromProtoToRegisterRequestDTO(req *authpb.RegisterRequest) *dto.RegisterRequest {
	return &dto.RegisterRequest{
		Name:      req.GetName(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
		AdminCode: req.GetAdminCode(),
	}
}

func ConvertFromProtoToLoginRequestDTO(req *authpb.LoginRequest) *dto.LoginRequest {
	return &dto.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func ConvertFromProtoToRefreshRequestDTO(req *authpb.RefreshRequest) *dto.RefreshRequest {
	return &dto.RefreshRequest{
		RefreshToken: req.GetRefreshToken(),
	}
}

func ConvertFromProtoToLogoutRequestDTO(req *authpb.LogoutRequest) *dto.LogoutRequest {
	return &dto.LogoutRequest{
		RefreshToken: req.GetRefreshToken(),
	}
}

func ConvertFromDTOToAuthResponse(d *dto.AuthResult) *authpb.AuthResponse {
	return &authpb.AuthResponse{
		UserId:       d.UserID,
		AccessToken:  d.AccessToken,
		RefreshToken: d.RefreshToken,
	}
}

func ConvertFromDTOToRefreshResponse(d *dto.RefreshResult) *authpb.RefreshResponse {
	return &authpb.RefreshResponse{
		AccessToken: d.AccessToken,
	}
}
