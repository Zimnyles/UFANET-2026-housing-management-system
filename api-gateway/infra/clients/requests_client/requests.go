package requests_client

import (
	"context"

	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/domain"
)

func (c *RequestsClient) CreateRequest(ctx context.Context, req *domain.CreateMaintenanceRequest) (*domain.MaintenanceRequest, error) {
	resp, err := c.client.CreateRequest(ctx, toProtoCreateRequest(req))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	return toDomainRequest(resp.GetRequest()), nil
}

func (c *RequestsClient) GetRequests(ctx context.Context, req *domain.ListMaintenanceRequests) ([]*domain.MaintenanceRequest, int64, error) {
	resp, err := c.client.GetRequests(ctx, toProtoGetRequests(req))
	if err != nil {
		return nil, 0, app_errors.FromGRPC(err)
	}

	result := make([]*domain.MaintenanceRequest, 0, len(resp.GetRequests()))
	for _, item := range resp.GetRequests() {
		result = append(result, toDomainRequest(item))
	}

	return result, int64(resp.GetTotal()), nil
}

func (c *RequestsClient) GetRequest(ctx context.Context, id string) (*domain.MaintenanceRequest, error) {
	resp, err := c.client.GetRequest(ctx, toProtoGetRequest(id))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	return toDomainRequest(resp.GetRequest()), nil
}

func (c *RequestsClient) UpdateRequestStatus(ctx context.Context, req *domain.UpdateMaintenanceRequestStatus) (*domain.MaintenanceRequest, error) {
	resp, err := c.client.UpdateRequestStatus(ctx, toProtoUpdateStatus(req))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	return toDomainRequest(resp.GetRequest()), nil
}

func (c *RequestsClient) AddComment(ctx context.Context, req *domain.AddMaintenanceRequestComment) (*domain.RequestComment, error) {
	resp, err := c.client.AddComment(ctx, toProtoAddComment(req))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	return toDomainComment(resp.GetComment()), nil
}

func (c *RequestsClient) GetComments(ctx context.Context, requestID string) ([]*domain.RequestComment, error) {
	resp, err := c.client.GetComments(ctx, toProtoGetComments(requestID))
	if err != nil {
		return nil, app_errors.FromGRPC(err)
	}

	result := make([]*domain.RequestComment, 0, len(resp.GetComments()))
	for _, item := range resp.GetComments() {
		result = append(result, toDomainComment(item))
	}

	return result, nil
}
