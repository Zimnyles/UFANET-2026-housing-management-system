package notifications_service

import (
	"context"

	"api-gateway/internal/models/domain"
)

type NotificationsClient interface {
	Register(context.Context, string, domain.RegisterDevice) error
	Unregister(context.Context, string, string) error
	Subscribe(context.Context, string, string) (<-chan domain.BrowserNotification, func(), error)
}
