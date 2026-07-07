package notifications_service

import (
	"context"
	"sync"

	"github.com/rs/zerolog"

	"notification-service/infra/models/domain"
)

type subscription struct {
	userID, houseID string
	events          chan *domain.Notification
}

type Service struct {
	repo   Repository
	logger *zerolog.Logger
	mu     sync.RWMutex
	nextID uint64
	subs   map[uint64]subscription
}

func New(repo Repository, logger *zerolog.Logger) *Service {
	return &Service{repo: repo, logger: logger, subs: make(map[uint64]subscription)}
}

func (s *Service) Register(ctx context.Context, device *domain.Device) error {
	return s.repo.Register(ctx, device)
}

func (s *Service) Unregister(ctx context.Context, userID, token string) error {
	return s.repo.Unregister(ctx, userID, token)
}

func (s *Service) List(ctx context.Context, req *domain.ListRequest) ([]*domain.Notification, int64, error) {
	return s.repo.List(ctx, req)
}

func (s *Service) Notify(ctx context.Context, n *domain.Notification) error {
	created, err := s.repo.Create(ctx, n)
	if err != nil {
		return err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, sub := range s.subs {
		if created.UserID != "" && created.UserID != sub.userID {
			continue
		}

		if created.HouseID != "" && created.HouseID != sub.houseID {
			continue
		}

		select {
		case sub.events <- created:
		default:
			s.logger.Warn().Str("user_id", sub.userID).Msg("notification subscriber is slow")
		}
	}

	return nil
}

func (s *Service) Subscribe(userID, houseID string) (<-chan *domain.Notification, func()) {
	s.mu.Lock()
	s.nextID++
	id := s.nextID
	events := make(chan *domain.Notification, 16)
	s.subs[id] = subscription{userID: userID, houseID: houseID, events: events}
	s.mu.Unlock()

	return events, func() {
		s.mu.Lock()
		if sub, ok := s.subs[id]; ok {
			delete(s.subs, id)
			close(sub.events)
		}
		s.mu.Unlock()
	}
}
