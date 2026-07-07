package notifications_handler

import (
	"context"

	"api-gateway/internal/models/domain"
)

type NotificationsService interface {
	Register(context.Context, string, domain.RegisterDevice) error
	Unregister(context.Context, string, string) error
	Subscribe(context.Context, string, string) (<-chan domain.BrowserNotification, func(), error)
}

type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (*domain.Profile, error)
}
