package requests_service

import (
	"context"

	"github.com/rs/zerolog"

	"api-gateway/internal/models/domain"
)

type RequestsService struct {
	client RequestsClient
	logger *zerolog.Logger
}

func New(client RequestsClient, logger *zerolog.Logger) *RequestsService {
	return &RequestsService{client: client, logger: logger}
}

func (s *RequestsService) CreateRequest(ctx context.Context, req *domain.CreateMaintenanceRequest) (*domain.MaintenanceRequest, error) {
	r, err := s.client.CreateRequest(ctx, req)
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", req.UserID).Msg("create request failed")

		return nil, err
	}

	return r, nil
}

func (s *RequestsService) GetRequests(ctx context.Context, req *domain.ListMaintenanceRequests) ([]*domain.MaintenanceRequest, int64, error) {
	return s.client.GetRequests(ctx, req)
}

func (s *RequestsService) GetRequest(ctx context.Context, id string) (*domain.MaintenanceRequest, error) {
	return s.client.GetRequest(ctx, id)
}

func (s *RequestsService) UpdateRequestStatus(ctx context.Context, req *domain.UpdateMaintenanceRequestStatus) (*domain.MaintenanceRequest, error) {
	r, err := s.client.UpdateRequestStatus(ctx, req)
	if err != nil {
		s.logger.Error().Err(err).Str("request_id", req.ID).Msg("update request status failed")

		return nil, err
	}

	return r, nil
}

func (s *RequestsService) AddComment(ctx context.Context, req *domain.AddMaintenanceRequestComment) (*domain.RequestComment, error) {
	c, err := s.client.AddComment(ctx, req)
	if err != nil {
		s.logger.Error().Err(err).Str("request_id", req.RequestID).Msg("add request comment failed")

		return nil, err
	}

	return c, nil
}

func (s *RequestsService) GetComments(ctx context.Context, requestID string) ([]*domain.RequestComment, error) {
	return s.client.GetComments(ctx, requestID)
}
