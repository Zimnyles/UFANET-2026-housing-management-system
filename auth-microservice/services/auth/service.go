package auth

import (
	infra_errors "auth-service/infra/errors"
	"auth-service/infra/models/domain"
	"auth-service/infra/models/dto"
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type AuthService struct {
	repo   Repository
	jwt    JWTManager
	hasher Hasher
	logger *zerolog.Logger
}

func New(repo Repository, jwt JWTManager, hasher Hasher, logger *zerolog.Logger) *AuthService {
	return &AuthService{
		repo:   repo,
		jwt:    jwt,
		hasher: hasher,
		logger: logger,
	}
}

func (s *AuthService) Register(ctx context.Context, d *dto.RegisterRequest) (*dto.AuthResult, error) {
	req := domain.ConvertFromDTOToRegisterRequest(d)

	passwordHash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, &domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         "user",
	})
	if err != nil {
		s.logger.Error().Err(err).Str("email", req.Email).Msg("register: create user failed")
		return nil, err
	}

	return s.issueTokens(user)
}

func (s *AuthService) Login(ctx context.Context, d *dto.LoginRequest) (*dto.AuthResult, error) {
	req := domain.ConvertFromDTOToLoginRequest(d)

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Error().Err(err).Str("email", req.Email).Msg("login: get user failed")
		return nil, err
	}

	if !s.hasher.Check(req.Password, user.PasswordHash) {
		return nil, infra_errors.ErrInvalidCredentials
	}

	return s.issueTokens(user)
}

func (s *AuthService) Refresh(ctx context.Context, d *dto.RefreshRequest) (*dto.RefreshResult, error) {
	req := domain.ConvertFromDTOToRefreshRequest(d)

	tokenClaims, err := s.jwt.ParseRefresh(req.RefreshToken)
	if err != nil {
		s.logger.Error().Err(err).Msg("refresh: invalid token")
		return nil, infra_errors.ErrInvalidToken
	}

	accessToken, err := s.jwt.GenerateAccess(tokenClaims.UserID, tokenClaims.Role)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	s.logger.Info().Str("user_id", tokenClaims.UserID).Msg("token refreshed")
	return &dto.RefreshResult{AccessToken: accessToken}, nil
}

func (s *AuthService) Logout(_ context.Context, _ *dto.LogoutRequest) error {
	return nil
}

func (s *AuthService) issueTokens(user *domain.User) (*dto.AuthResult, error) {
	accessToken, err := s.jwt.GenerateAccess(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.jwt.GenerateRefresh(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	s.logger.Info().Str("user_id", user.ID).Msg("tokens issued")
	return &dto.AuthResult{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
