package news_service

import (
	"context"
	"fmt"
	"news-service/infra/models/domain"

	"github.com/rs/zerolog"
)

type NewsService struct {
	repo      Repository
	publisher Publisher
	logger    *zerolog.Logger
}

func New(repo Repository, publisher Publisher, logger *zerolog.Logger) *NewsService {
	return &NewsService{repo: repo, publisher: publisher, logger: logger}
}

func (s *NewsService) CreateNews(ctx context.Context, req *domain.CreateNewsRequest) (*domain.News, error) {
	n, err := s.repo.Create(ctx, &domain.News{
		Title:   req.Title,
		Content: req.Content,
		HouseID: req.HouseID,
	})
	if err != nil {
		s.logger.Error().Err(err).Str("house_id", req.HouseID).Msg("create news failed")
		return nil, fmt.Errorf("create news: %w", err)
	}

	if s.publisher != nil {
		evt := domain.NewsCreatedEvent{
			Type:    "news.created",
			NewsID:  n.ID,
			Title:   n.Title,
			HouseID: n.HouseID,
		}
		if pErr := s.publisher.PublishNewsCreated(ctx, evt); pErr != nil {
			s.logger.Warn().Err(pErr).Str("news_id", n.ID).Msg("publish news.created failed (news saved)")
		}
	}

	s.logger.Info().Str("id", n.ID).Str("house_id", n.HouseID).Msg("news created")
	return n, nil
}

func (s *NewsService) GetNewsItem(ctx context.Context, id string) (*domain.News, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *NewsService) GetNews(ctx context.Context, req *domain.GetNewsRequest) ([]*domain.News, int64, error) {
	return s.repo.List(ctx, req)
}
