package news_service

import (
	"context"

	"github.com/rs/zerolog"
)

type NewsService struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *NewsService {
	return &NewsService{logger: logger}
}

func (s *NewsService) List(ctx context.Context) error {
	return nil
}

func (s *NewsService) Create(ctx context.Context) error {
	return nil
}
