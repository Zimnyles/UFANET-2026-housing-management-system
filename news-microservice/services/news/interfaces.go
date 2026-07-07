package news_service

import (
	"context"

	"news-service/infra/models/domain"
)

type Repository interface {
	Create(ctx context.Context, n *domain.News) (*domain.News, error)
	GetByID(ctx context.Context, id string) (*domain.News, error)
	List(ctx context.Context, req *domain.GetNewsRequest) ([]*domain.News, int64, error)
}

type Publisher interface {
	PublishNewsCreated(ctx context.Context, evt domain.NewsCreatedEvent) error
}
