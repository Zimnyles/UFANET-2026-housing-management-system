package server

import (
	"context"
	"math"

	"github.com/rs/zerolog"
	requestspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/requests/langs/go"
)

type RequestsServer struct {
	requestspb.UnimplementedRequestsServiceServer
	service RequestsService
	logger  *zerolog.Logger
}

func NewRequestsServer(service RequestsService, logger *zerolog.Logger) *RequestsServer {
	return &RequestsServer{service: service, logger: logger}
}

func (s *RequestsServer) CreateRequest(ctx context.Context, req *requestspb.CreateRequestRequest) (*requestspb.RequestResponse, error) {
	r, err := s.service.CreateRequest(ctx, protoToCreateRequest(req))
	if err != nil {
		return nil, err
	}

	return &requestspb.RequestResponse{Request: domainToProtoRequest(r)}, nil
}

func (s *RequestsServer) GetRequests(ctx context.Context, req *requestspb.GetRequestsRequest) (*requestspb.GetRequestsResponse, error) {
	items, total, err := s.service.GetRequests(ctx, protoToGetRequests(req))
	if err != nil {
		return nil, err
	}

	pb := make([]*requestspb.MaintenanceRequest, 0, len(items))
	for _, item := range items {
		pb = append(pb, domainToProtoRequest(item))
	}

	return &requestspb.GetRequestsResponse{Requests: pb, Total: totalToInt32(total)}, nil
}

func totalToInt32(total int64) int32 {
	if total <= 0 {
		return 0
	}

	if total > math.MaxInt32 {
		return math.MaxInt32
	}

	return int32(total)
}

func (s *RequestsServer) GetRequest(ctx context.Context, req *requestspb.GetRequestRequest) (*requestspb.RequestResponse, error) {
	r, err := s.service.GetRequest(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &requestspb.RequestResponse{Request: domainToProtoRequest(r)}, nil
}

func (s *RequestsServer) UpdateRequestStatus(ctx context.Context, req *requestspb.UpdateStatusRequest) (*requestspb.RequestResponse, error) {
	r, err := s.service.UpdateRequestStatus(ctx, protoToUpdateStatus(req))
	if err != nil {
		return nil, err
	}

	return &requestspb.RequestResponse{Request: domainToProtoRequest(r)}, nil
}

func (s *RequestsServer) AddComment(ctx context.Context, req *requestspb.AddCommentRequest) (*requestspb.CommentResponse, error) {
	c, err := s.service.AddComment(ctx, protoToAddComment(req))
	if err != nil {
		return nil, err
	}

	return &requestspb.CommentResponse{Comment: domainToProtoComment(c)}, nil
}

func (s *RequestsServer) GetComments(ctx context.Context, req *requestspb.GetCommentsRequest) (*requestspb.GetCommentsResponse, error) {
	items, err := s.service.GetComments(ctx, req.GetRequestId())
	if err != nil {
		return nil, err
	}

	pb := make([]*requestspb.Comment, 0, len(items))
	for _, item := range items {
		pb = append(pb, domainToProtoComment(item))
	}

	return &requestspb.GetCommentsResponse{Comments: pb}, nil
}
