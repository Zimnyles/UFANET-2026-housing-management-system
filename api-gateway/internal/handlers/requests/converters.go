package requests_handler

import (
	"api-gateway/internal/models/domain"
	"api-gateway/internal/models/dto"
)

func toDomainCreateRequest(userID string, req dto.CreateMaintenanceRequest) *domain.CreateMaintenanceRequest {
	return &domain.CreateMaintenanceRequest{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		UserID:      userID,
	}
}

func toDTORequest(r *domain.MaintenanceRequest) dto.MaintenanceRequestResponse {
	return dto.MaintenanceRequestResponse{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		Type:        r.Type,
		Status:      r.Status,
		UserID:      r.UserID,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toDTOComment(c *domain.RequestComment) dto.RequestCommentResponse {
	return dto.RequestCommentResponse{
		ID:        c.ID,
		RequestID: c.RequestID,
		UserID:    c.UserID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}
