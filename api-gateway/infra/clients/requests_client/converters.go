package requests_client

import (
	"math"

	requestspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/requests/langs/go"

	"api-gateway/internal/models/domain"
)

func toProtoCreateRequest(req *domain.CreateMaintenanceRequest) *requestspb.CreateRequestRequest {
	return &requestspb.CreateRequestRequest{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		UserId:      req.UserID,
	}
}

func toProtoGetRequests(req *domain.ListMaintenanceRequests) *requestspb.GetRequestsRequest {
	return &requestspb.GetRequestsRequest{
		UserId: req.UserID,
		Status: req.Status,
		Type:   req.Type,
		Limit:  paginationInt32(req.Limit),
		Offset: paginationInt32(req.Offset),
	}
}

func paginationInt32(value int) int32 {
	if value <= 0 {
		return 0
	}

	if value > math.MaxInt32 {
		return math.MaxInt32
	}

	return int32(value)
}

func toProtoGetRequest(id string) *requestspb.GetRequestRequest {
	return &requestspb.GetRequestRequest{Id: id}
}

func toProtoUpdateStatus(req *domain.UpdateMaintenanceRequestStatus) *requestspb.UpdateStatusRequest {
	return &requestspb.UpdateStatusRequest{
		Id:     req.ID,
		Status: req.Status,
		UserId: req.UserID,
	}
}

func toProtoAddComment(req *domain.AddMaintenanceRequestComment) *requestspb.AddCommentRequest {
	return &requestspb.AddCommentRequest{
		RequestId: req.RequestID,
		UserId:    req.UserID,
		Content:   req.Content,
	}
}

func toProtoGetComments(requestID string) *requestspb.GetCommentsRequest {
	return &requestspb.GetCommentsRequest{RequestId: requestID}
}

func toDomainRequest(pb *requestspb.MaintenanceRequest) *domain.MaintenanceRequest {
	if pb == nil {
		return nil
	}

	return &domain.MaintenanceRequest{
		ID:          pb.GetId(),
		Title:       pb.GetTitle(),
		Description: pb.GetDescription(),
		Type:        pb.GetType(),
		Status:      pb.GetStatus(),
		UserID:      pb.GetUserId(),
		CreatedAt:   pb.GetCreatedAt(),
		UpdatedAt:   pb.GetUpdatedAt(),
	}
}

func toDomainComment(pb *requestspb.Comment) *domain.RequestComment {
	if pb == nil {
		return nil
	}

	return &domain.RequestComment{
		ID:        pb.GetId(),
		RequestID: pb.GetRequestId(),
		UserID:    pb.GetUserId(),
		Content:   pb.GetContent(),
		CreatedAt: pb.GetCreatedAt(),
	}
}
