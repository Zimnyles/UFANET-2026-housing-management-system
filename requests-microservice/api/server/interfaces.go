package server

import (
	"context"

	"requests-service/infra/models/domain"
)

type RequestsService interface {
	CreateRequest(ctx context.Context, req *domain.CreateRequestRequest) (*domain.MaintenanceRequest, error)
	GetRequests(ctx context.Context, req *domain.GetRequestsRequest) ([]*domain.MaintenanceRequest, int64, error)
	GetRequest(ctx context.Context, id string) (*domain.MaintenanceRequest, error)
	UpdateRequestStatus(ctx context.Context, req *domain.UpdateStatusRequest) (*domain.MaintenanceRequest, error)
	AddComment(ctx context.Context, req *domain.AddCommentRequest) (*domain.Comment, error)
	GetComments(ctx context.Context, requestID string) ([]*domain.Comment, error)
}
