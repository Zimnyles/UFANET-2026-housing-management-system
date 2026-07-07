package server

import (
	"context"

	"notification-service/infra/models/domain"
)

type NotificationsService interface {
	Register(context.Context, *domain.Device) error
	Unregister(context.Context, string, string) error
	List(context.Context, *domain.ListRequest) ([]*domain.Notification, int64, error)
	Subscribe(string, string) (<-chan *domain.Notification, func())
}
