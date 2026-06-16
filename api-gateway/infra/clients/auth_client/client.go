package auth_client

import (
	"fmt"

	app_errors "api-gateway/internal/errors"
	"github.com/rs/zerolog"
	authpb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/auth/langs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn   ClientConnection
	client AuthServiceClient
	logger *zerolog.Logger
}

func New(addr string, logger *zerolog.Logger) (*AuthClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: auth at %s: %w", app_errors.ErrClientDial, addr, err)
	}

	logger.Info().Str("addr", addr).Msg("auth client connected")

	return &AuthClient{
		conn:   conn,
		client: authpb.NewAuthServiceClient(conn),
		logger: logger,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}
