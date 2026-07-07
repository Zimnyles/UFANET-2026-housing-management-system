package news_client

import (
	"context"

	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/domain"
)

func (c *Client) Create(ctx context.Context, req *domain.CreateNews) (*domain.News, error) {
	resp, err := c.client.CreateNews(ctx, toProtoCreate(req))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	return toDomain(resp.GetNews()), nil
}

func (c *Client) List(ctx context.Context, req *domain.ListNews) ([]*domain.News, int64, error) {
	resp, err := c.client.GetNews(ctx, toProtoList(req))
	if err != nil {
		return nil, 0, app_errors.FromGRPC(err)
	}

	return toDomainList(resp.GetNews()), int64(resp.GetTotal()), nil
}
