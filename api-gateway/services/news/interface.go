package news_service

import (
	"context"

	"api-gateway/internal/models/domain"
)

type NewsClient interface {
	Create(context.Context, *domain.CreateNews) (*domain.News, error)
	List(context.Context, *domain.ListNews) ([]*domain.News, int64, error)
}
