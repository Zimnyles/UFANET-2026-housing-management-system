package auth_client

import (
	"context"

	"api-gateway/internal/models/domain"
)

func (c *AuthClient) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResult, error) {
	pb, err := c.client.Register(ctx, toProtoRegisterRequest(req))
	if err != nil {
		return nil, err
	}

	return toDomainAuthResult(pb), nil
}

func (c *AuthClient) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResult, error) {
	pb, err := c.client.Login(ctx, toProtoLoginRequest(req))
	if err != nil {
		return nil, err
	}

	return toDomainAuthResult(pb), nil
}

func (c *AuthClient) Refresh(ctx context.Context, req *domain.RefreshRequest) (domain.TokenPair, error) {
	pb, err := c.client.Refresh(ctx, toProtoRefreshRequest(req))
	if err != nil {
		return domain.TokenPair{}, err
	}

	return toDomainTokenPair(pb), nil
}

func (c *AuthClient) Logout(ctx context.Context, refreshToken string) error {
	_, err := c.client.Logout(ctx, toProtoLogoutRequest(refreshToken))

	return err
}
