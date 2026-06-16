package profile_client

import (
	"fmt"

	app_errors "api-gateway/internal/errors"

	"github.com/rs/zerolog"
	profilepb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/profile/langs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProfileClient struct {
	conn   ClientConnection
	client ProfileServiceClient
	logger *zerolog.Logger
}

func New(addr string, logger *zerolog.Logger) (*ProfileClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: profile at %s: %w", app_errors.ErrClientDial, addr, err)
	}

	logger.Info().Str("addr", addr).Msg("profile client connected")

	return &ProfileClient{
		conn:   conn,
		client: profilepb.NewProfileServiceClient(conn),
		logger: logger,
	}, nil
}

func (c *ProfileClient) Close() error {
	return c.conn.Close()
}
