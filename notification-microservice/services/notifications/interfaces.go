package notifications_service

import (
	"context"

	"notification-service/infra/models/domain"
)

type Repository interface {
	Register(context.Context, *domain.Device) error
	Unregister(context.Context, string, string) error
	Create(context.Context, *domain.Notification) (*domain.Notification, error)
	List(context.Context, *domain.ListRequest) ([]*domain.Notification, int64, error)
}
