package server

import (
	"auth-service/infra/models/dto"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, d *dto.RegisterRequest) (*dto.AuthResult, error)
	Login(ctx context.Context, d *dto.LoginRequest) (*dto.AuthResult, error)
	Refresh(ctx context.Context, d *dto.RefreshRequest) (*dto.RefreshResult, error)
	Logout(ctx context.Context, d *dto.LogoutRequest) error
}
