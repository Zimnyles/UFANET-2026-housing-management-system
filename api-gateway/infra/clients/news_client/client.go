package news_client

import (
	"fmt"

	"github.com/rs/zerolog"
	newspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/news/langs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	app_errors "api-gateway/internal/errors"
)

type Client struct {
	conn   ClientConnection
	client NewsServiceClient
	logger *zerolog.Logger
}

func New(addr string, logger *zerolog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%w: news at %s: %w", app_errors.ErrClientDial, addr, err)
	}

	return &Client{conn: conn, client: newspb.NewNewsServiceClient(conn), logger: logger}, nil
}

func (c *Client) Close() error { return c.conn.Close() }
