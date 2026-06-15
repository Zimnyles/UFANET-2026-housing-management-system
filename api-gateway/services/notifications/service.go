package notifications_service

import (
	"context"

	"github.com/rs/zerolog"
)

type NotificationsService struct {
	logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *NotificationsService {
	return &NotificationsService{logger: logger}
}

func (s *NotificationsService) Register(ctx context.Context) error {
	return nil
}

func (s *NotificationsService) Unregister(ctx context.Context) error {
	return nil
}
