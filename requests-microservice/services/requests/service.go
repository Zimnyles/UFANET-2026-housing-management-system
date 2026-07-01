package requests_service

import (
	"context"
	"fmt"
	"requests-service/infra/models/domain"

	"github.com/rs/zerolog"
)

type RequestsService struct {
	repo      Repository
	publisher Publisher
	logger    *zerolog.Logger
}

func New(repo Repository, publisher Publisher, logger *zerolog.Logger) *RequestsService {
	return &RequestsService{repo: repo, publisher: publisher, logger: logger}
}

func (s *RequestsService) CreateRequest(ctx context.Context, req *domain.CreateRequestRequest) (*domain.MaintenanceRequest, error) {
	r, err := s.repo.Create(ctx, &domain.MaintenanceRequest{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		UserID:      req.UserID,
	})
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", req.UserID).Msg("create request failed")
		return nil, fmt.Errorf("create request: %w", err)
	}

	s.publishCreated(ctx, r)
	s.logger.Info().Str("id", r.ID).Str("user_id", r.UserID).Msg("request created")
	return r, nil
}

func (s *RequestsService) GetRequests(ctx context.Context, req *domain.GetRequestsRequest) ([]*domain.MaintenanceRequest, int64, error) {
	return s.repo.List(ctx, req)
}

func (s *RequestsService) GetRequest(ctx context.Context, id string) (*domain.MaintenanceRequest, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RequestsService) UpdateRequestStatus(ctx context.Context, req *domain.UpdateStatusRequest) (*domain.MaintenanceRequest, error) {
	r, err := s.repo.UpdateStatus(ctx, req.ID, req.Status)
	if err != nil {
		s.logger.Error().Err(err).Str("request_id", req.ID).Msg("update request status failed")
		return nil, fmt.Errorf("update request status: %w", err)
	}

	s.publishStatusUpdated(ctx, r, req.UserID)
	s.logger.Info().Str("id", r.ID).Str("status", r.Status).Msg("request status updated")
	return r, nil
}

func (s *RequestsService) AddComment(ctx context.Context, req *domain.AddCommentRequest) (*domain.Comment, error) {
	comment, err := s.repo.AddComment(ctx, &domain.Comment{
		RequestID: req.RequestID,
		UserID:    req.UserID,
		Content:   req.Content,
	})
	if err != nil {
		s.logger.Error().Err(err).Str("request_id", req.RequestID).Msg("add comment failed")
		return nil, fmt.Errorf("add comment: %w", err)
	}

	s.publishCommentAdded(ctx, comment)
	s.logger.Info().Str("id", comment.ID).Str("request_id", comment.RequestID).Msg("comment added")
	return comment, nil
}

func (s *RequestsService) GetComments(ctx context.Context, requestID string) ([]*domain.Comment, error) {
	return s.repo.ListComments(ctx, requestID)
}

func (s *RequestsService) publishCreated(ctx context.Context, r *domain.MaintenanceRequest) {
	if s.publisher == nil {
		return
	}
	evt := domain.RequestEvent{
		Type:      "request.created",
		RequestID: r.ID,
		UserID:    r.UserID,
		Title:     r.Title,
		Status:    r.Status,
	}
	if err := s.publisher.PublishRequestCreated(ctx, evt); err != nil {
		s.logger.Warn().Err(err).Str("request_id", r.ID).Msg("publish request.created failed")
	}
}

func (s *RequestsService) publishStatusUpdated(ctx context.Context, r *domain.MaintenanceRequest, actorID string) {
	if s.publisher == nil {
		return
	}
	evt := domain.RequestEvent{
		Type:      "request.status_updated",
		RequestID: r.ID,
		UserID:    actorID,
		Status:    r.Status,
	}
	if err := s.publisher.PublishRequestStatusUpdated(ctx, evt); err != nil {
		s.logger.Warn().Err(err).Str("request_id", r.ID).Msg("publish request.status_updated failed")
	}
}

func (s *RequestsService) publishCommentAdded(ctx context.Context, comment *domain.Comment) {
	if s.publisher == nil {
		return
	}
	evt := domain.RequestEvent{
		Type:      "request.comment_added",
		RequestID: comment.RequestID,
		UserID:    comment.UserID,
	}
	if err := s.publisher.PublishRequestCommentAdded(ctx, evt); err != nil {
		s.logger.Warn().Err(err).Str("request_id", comment.RequestID).Msg("publish request.comment_added failed")
	}
}
