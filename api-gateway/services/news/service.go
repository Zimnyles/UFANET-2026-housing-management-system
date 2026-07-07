package news_service

import (
	"context"

	"github.com/rs/zerolog"

	"api-gateway/internal/models/domain"
)

type NewsService struct {
	client NewsClient
	logger *zerolog.Logger
}

func New(client NewsClient, logger *zerolog.Logger) *NewsService {
	return &NewsService{client: client, logger: logger}
}

func (s *NewsService) List(ctx context.Context, req *domain.ListNews) ([]*domain.News, int64, error) {
	return s.client.List(ctx, req)
}

func (s *NewsService) Create(ctx context.Context, req *domain.CreateNews) (*domain.News, error) {
	return s.client.Create(ctx, req)
}
