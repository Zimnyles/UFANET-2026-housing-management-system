package auth_service

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/domain"
	"context"

	"github.com/rs/zerolog"
)

type AuthService struct {
	authClient AuthClient
	logger     *zerolog.Logger
}

func New(authClient AuthClient, logger *zerolog.Logger) *AuthService {
	return &AuthService{authClient: authClient, logger: logger}
}

func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResult, error) {
	result, err := s.authClient.Register(ctx, req)
	if err != nil {
		s.logger.Error().Err(err).Str("email", req.Email).Msg("register failed")
		return nil, app_errors.FromGRPC(err)
	}
	s.logger.Info().Str("user_id", result.UserID).Msg("registered")
	return result, nil
}

func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResult, error) {
	result, err := s.authClient.Login(ctx, req)
	if err != nil {
		s.logger.Error().Err(err).Str("email", req.Email).Msg("login failed")
		return nil, app_errors.FromGRPC(err)
	}
	s.logger.Info().Str("user_id", result.UserID).Msg("logged in")
	return result, nil
}

func (s *AuthService) Refresh(ctx context.Context, req *domain.RefreshRequest) (*domain.TokenPair, error) {
	tokens, err := s.authClient.Refresh(ctx, req)
	if err != nil {
		s.logger.Error().Err(err).Msg("refresh failed")
		return nil, app_errors.FromGRPC(err)
	}
	return &tokens, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.authClient.Logout(ctx, refreshToken); err != nil {
		s.logger.Error().Err(err).Msg("logout failed")
		return app_errors.FromGRPC(err)
	}
	return nil
}
