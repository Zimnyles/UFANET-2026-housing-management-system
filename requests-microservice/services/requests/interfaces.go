package requests_service

import (
	"context"

	"requests-service/infra/models/domain"
)

type Repository interface {
	Create(ctx context.Context, req *domain.MaintenanceRequest) (*domain.MaintenanceRequest, error)
	GetByID(ctx context.Context, id string) (*domain.MaintenanceRequest, error)
	List(ctx context.Context, req *domain.GetRequestsRequest) ([]*domain.MaintenanceRequest, int64, error)
	UpdateStatus(ctx context.Context, id string, status string) (*domain.MaintenanceRequest, error)
	AddComment(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	ListComments(ctx context.Context, requestID string) ([]*domain.Comment, error)
}

type Publisher interface {
	PublishRequestCreated(ctx context.Context, evt domain.RequestEvent) error
	PublishRequestStatusUpdated(ctx context.Context, evt domain.RequestEvent) error
	PublishRequestCommentAdded(ctx context.Context, evt domain.RequestEvent) error
}
