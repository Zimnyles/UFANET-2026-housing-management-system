package requests_client

import (
	"fmt"

	app_errors "api-gateway/internal/errors"

	"github.com/rs/zerolog"
	requestspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/requests/langs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RequestsClient struct {
	conn   ClientConnection
	client RequestsServiceClient
	logger *zerolog.Logger
}

func New(addr string, logger *zerolog.Logger) (*RequestsClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: requests at %s: %w", app_errors.ErrClientDial, addr, err)
	}

	logger.Info().Str("addr", addr).Msg("requests client connected")

	return &RequestsClient{
		conn:   conn,
		client: requestspb.NewRequestsServiceClient(conn),
		logger: logger,
	}, nil
}

func (c *RequestsClient) Close() error {
	return c.conn.Close()
}
