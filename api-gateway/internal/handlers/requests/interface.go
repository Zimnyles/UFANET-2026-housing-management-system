package requests_handler

import (
	"context"

	"api-gateway/internal/models/domain"
)

type RequestsService interface {
	CreateRequest(ctx context.Context, req *domain.CreateMaintenanceRequest) (*domain.MaintenanceRequest, error)
	GetRequests(ctx context.Context, req *domain.ListMaintenanceRequests) ([]*domain.MaintenanceRequest, int64, error)
	GetRequest(ctx context.Context, id string) (*domain.MaintenanceRequest, error)
	UpdateRequestStatus(ctx context.Context, req *domain.UpdateMaintenanceRequestStatus) (*domain.MaintenanceRequest, error)
	AddComment(ctx context.Context, req *domain.AddMaintenanceRequestComment) (*domain.RequestComment, error)
	GetComments(ctx context.Context, requestID string) ([]*domain.RequestComment, error)
}

type ProfileService interface {
	GetProfile(context.Context, string) (*domain.Profile, error)
}
