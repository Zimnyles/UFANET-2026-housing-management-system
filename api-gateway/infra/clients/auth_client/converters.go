package auth_client

import (
	"api-gateway/internal/models/domain"

	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"
)

func toProtoRegisterRequest(req *domain.RegisterRequest) *authpb.RegisterRequest {
	return &authpb.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		AdminCode: req.AdminCode,
	}
}

func toProtoLoginRequest(req *domain.LoginRequest) *authpb.LoginRequest {
	return &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func toProtoRefreshRequest(req *domain.RefreshRequest) *authpb.RefreshRequest {
	return &authpb.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}
}

func toProtoLogoutRequest(refreshToken string) *authpb.LogoutRequest {
	return &authpb.LogoutRequest{
		RefreshToken: refreshToken,
	}
}

func toDomainAuthResult(pb *authpb.AuthResponse) *domain.AuthResult {
	return &domain.AuthResult{
		UserID: pb.GetUserId(),
		Tokens: domain.TokenPair{
			AccessToken:  pb.GetAccessToken(),
			RefreshToken: pb.GetRefreshToken(),
		},
	}
}

func toDomainTokenPair(pb *authpb.RefreshResponse) domain.TokenPair {
	return domain.TokenPair{
		AccessToken: pb.GetAccessToken(),
	}
}
