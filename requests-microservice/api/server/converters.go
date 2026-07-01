package server

import (
	"time"

	"requests-service/infra/models/domain"

	requestspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/requests/langs/go"
)

func protoToCreateRequest(req *requestspb.CreateRequestRequest) *domain.CreateRequestRequest {
	return &domain.CreateRequestRequest{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Type:        req.GetType(),
		UserID:      req.GetUserId(),
	}
}

func protoToGetRequests(req *requestspb.GetRequestsRequest) *domain.GetRequestsRequest {
	return &domain.GetRequestsRequest{
		UserID: req.GetUserId(),
		Status: req.GetStatus(),
		Type:   req.GetType(),
		Limit:  int(req.GetLimit()),
		Offset: int(req.GetOffset()),
	}
}

func protoToUpdateStatus(req *requestspb.UpdateStatusRequest) *domain.UpdateStatusRequest {
	return &domain.UpdateStatusRequest{
		ID:     req.GetId(),
		Status: req.GetStatus(),
		UserID: req.GetUserId(),
	}
}

func protoToAddComment(req *requestspb.AddCommentRequest) *domain.AddCommentRequest {
	return &domain.AddCommentRequest{
		RequestID: req.GetRequestId(),
		UserID:    req.GetUserId(),
		Content:   req.GetContent(),
	}
}

func domainToProtoRequest(r *domain.MaintenanceRequest) *requestspb.MaintenanceRequest {
	return &requestspb.MaintenanceRequest{
		Id:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		Type:        r.Type,
		Status:      r.Status,
		UserId:      r.UserID,
		CreatedAt:   r.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func domainToProtoComment(c *domain.Comment) *requestspb.Comment {
	return &requestspb.Comment{
		Id:        c.ID,
		RequestId: c.RequestID,
		UserId:    c.UserID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt.UTC().Format(time.RFC3339),
	}
}
