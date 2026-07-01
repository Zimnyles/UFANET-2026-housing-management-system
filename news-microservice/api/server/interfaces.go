package server

import (
	"context"
	"news-service/infra/models/domain"
)

type NewsService interface {
	CreateNews(ctx context.Context, req *domain.CreateNewsRequest) (*domain.News, error)
	GetNewsItem(ctx context.Context, id string) (*domain.News, error)
	GetNews(ctx context.Context, req *domain.GetNewsRequest) ([]*domain.News, int64, error)
}
