package auth_service

import (
	"context"

	"api-gateway/internal/models/domain"
)

type AuthClient interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResult, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResult, error)
	Refresh(ctx context.Context, req *domain.RefreshRequest) (domain.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}
