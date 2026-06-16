package auth_handler

import (
	"context"
	"time"

	"api-gateway/internal/models/domain"
)

type AuthService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResult, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResult, error)
	Refresh(ctx context.Context, req *domain.RefreshRequest) (*domain.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
}

type Blacklist interface {
	BlacklistToken(token string, ttl time.Duration) error
	BlacklistRawJWT(token string) error
	IsBlacklisted(token string) bool
}
