package requests_service

import "requests-service/infra/models/domain"

func createRequestToModel(req *domain.CreateRequestRequest) *domain.MaintenanceRequest {
	return &domain.MaintenanceRequest{
		Title: req.Title, Description: req.Description, Type: req.Type, UserID: req.UserID,
	}
}

func addCommentToModel(req *domain.AddCommentRequest) *domain.Comment {
	return &domain.Comment{RequestID: req.RequestID, UserID: req.UserID, Content: req.Content}
}

func requestCreatedEvent(r *domain.MaintenanceRequest) domain.RequestEvent {
	return domain.RequestEvent{Type: "request.created", RequestID: r.ID, UserID: r.UserID, Title: r.Title, Status: r.Status}
}

func requestStatusUpdatedEvent(r *domain.MaintenanceRequest) domain.RequestEvent {
	return domain.RequestEvent{Type: "request.status_updated", RequestID: r.ID, UserID: r.UserID, Status: r.Status, Title: r.Title}
}

func requestCommentAddedEvent(comment *domain.Comment) domain.RequestEvent {
	return domain.RequestEvent{Type: "request.comment_added", RequestID: comment.RequestID, UserID: comment.UserID}
}
