package notifications_client

import (
	"context"
	"time"

	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/domain"
)

func (c *Client) Register(ctx context.Context, userID string, req domain.RegisterDevice) error {
	_, err := c.client.RegisterDevice(ctx, toProtoRegister(userID, req))
	if err != nil {
		return app_errors.FromGRPC(err)
	}

	return nil
}

func (c *Client) Unregister(ctx context.Context, userID, token string) error {
	_, err := c.client.UnregisterDevice(ctx, toProtoUnregister(userID, token))
	if err != nil {
		return app_errors.FromGRPC(err)
	}

	return nil
}

func (c *Client) Subscribe(ctx context.Context, userID, houseID string) (<-chan domain.BrowserNotification, func(), error) {
	streamCtx, cancel := context.WithCancel(ctx)

	stream, err := c.client.Subscribe(streamCtx, toProtoSubscribe(userID, houseID))
	if err != nil {
		cancel()

		return nil, nil, app_errors.FromGRPC(err)
	}

	events := make(chan domain.BrowserNotification, 16)

	go func() {
		defer close(events)

		for {
			item, err := stream.Recv()
			if err != nil {
				return
			}

			select {
			case events <- toDomainNotification(item):
			case <-streamCtx.Done():
				return
			}
		}
	}()

	return events, cancel, nil
}

func parseTime(value string) time.Time {
	parsed, _ := time.Parse(time.RFC3339, value)

	return parsed
}
