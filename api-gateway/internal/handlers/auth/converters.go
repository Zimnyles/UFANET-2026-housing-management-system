package auth_handler

import (
	"api-gateway/internal/models/domain"
	"api-gateway/internal/models/dto"
)

func toDomainRegisterRequest(req dto.RegisterRequest) *domain.RegisterRequest {
	return &domain.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func toDomainLoginRequest(req dto.LoginRequest) *domain.LoginRequest {
	return &domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}

func toDomainRefreshRequest(req dto.RefreshRequest) *domain.RefreshRequest {
	return &domain.RefreshRequest{
		RefreshToken: req.RefreshToken,
	}
}

func toDTORegisterResponse(res *domain.AuthResult) dto.RegisterResponse {
	return dto.RegisterResponse{
		UserID:       res.UserID,
		AccessToken:  res.Tokens.AccessToken,
		RefreshToken: res.Tokens.RefreshToken,
	}
}

func toDTOLoginResponse(res *domain.AuthResult) dto.LoginResponse {
	return dto.LoginResponse{
		UserID:       res.UserID,
		AccessToken:  res.Tokens.AccessToken,
		RefreshToken: res.Tokens.RefreshToken,
	}
}

func toDTORefreshResponse(tokens *domain.TokenPair) dto.RefreshResponse {
	return dto.RefreshResponse{
		AccessToken: tokens.AccessToken,
	}
}
