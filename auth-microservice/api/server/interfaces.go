package server

import (
	"auth-service/infra/models/domain"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, d *domain.RegisterRequest) (*domain.AuthResult, error)
	Login(ctx context.Context, d *domain.LoginRequest) (*domain.AuthResult, error)
	Refresh(ctx context.Context, d *domain.RefreshRequest) (*domain.RefreshResult, error)
	Logout(ctx context.Context, d *domain.LogoutRequest) error
}
