package notifications_service

import (
	"context"

	"github.com/rs/zerolog"

	"api-gateway/internal/models/domain"
)

type Service struct {
	client NotificationsClient
	logger *zerolog.Logger
}

func New(client NotificationsClient, logger *zerolog.Logger) *Service {
	return &Service{client: client, logger: logger}
}

func (s *Service) Register(ctx context.Context, userID string, device domain.RegisterDevice) error {
	if err := s.client.Register(ctx, userID, device); err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("register notification device failed")

		return err
	}

	return nil
}

func (s *Service) Unregister(ctx context.Context, userID, deviceToken string) error {
	if err := s.client.Unregister(ctx, userID, deviceToken); err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("unregister notification device failed")

		return err
	}

	return nil
}

func (s *Service) Subscribe(ctx context.Context, userID, houseID string) (<-chan domain.BrowserNotification, func(), error) {
	events, cancel, err := s.client.Subscribe(ctx, userID, houseID)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("subscribe to notifications failed")

		return nil, nil, err
	}

	return events, cancel, nil
}
