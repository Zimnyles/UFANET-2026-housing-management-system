package notifications_client

import (
	"fmt"

	notificationspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/notifications/langs/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	app_errors "api-gateway/internal/errors"
)

type Client struct {
	conn   ClientConnection
	client NotificationServiceClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%w: notifications at %s: %w", app_errors.ErrClientDial, addr, err)
	}

	return &Client{conn: conn, client: notificationspb.NewNotificationServiceClient(conn)}, nil
}
func (c *Client) Close() error { return c.conn.Close() }
