package news_handler

import (
	"context"

	"api-gateway/internal/models/domain"
)

type NewsService interface {
	Create(context.Context, *domain.CreateNews) (*domain.News, error)
	List(context.Context, *domain.ListNews) ([]*domain.News, int64, error)
}

type ProfileService interface {
	GetProfile(context.Context, string) (*domain.Profile, error)
	ListHouses(context.Context, string) ([]*domain.House, error)
}
