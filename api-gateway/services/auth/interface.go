package auth_service

import (
	"api-gateway/internal/models/domain"
	"context"
)

type AuthClient interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResult, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResult, error)
	Refresh(ctx context.Context, req *domain.RefreshRequest) (domain.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}
