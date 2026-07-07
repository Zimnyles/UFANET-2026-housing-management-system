package requests_handler

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/dto"
)

func validateCreateRequest(req *dto.CreateMaintenanceRequest) error {
	if req.Title == "" {
		return app_errors.ErrRequestTitleRequired
	}

	if len(req.Title) < 3 {
		return app_errors.ErrRequestTitleTooShort
	}

	if len(req.Title) > 255 {
		return app_errors.ErrRequestTitleTooLong
	}

	if req.Description == "" {
		return app_errors.ErrRequestDescriptionRequired
	}

	if len(req.Description) < 10 {
		return app_errors.ErrRequestDescriptionTooShort
	}

	if req.Type == "" {
		return app_errors.ErrRequestTypeRequired
	}

	if req.Type != "plumber" && req.Type != "electrician" {
		return app_errors.ErrRequestTypeInvalid
	}

	return nil
}

func validateUpdateStatus(req *dto.UpdateMaintenanceRequestStatus) error {
	if req.Status == "" {
		return app_errors.ErrStatusRequired
	}

	if req.Status != "open" && req.Status != "in_progress" && req.Status != "done" && req.Status != "cancelled" {
		return app_errors.ErrStatusInvalid
	}

	return nil
}

func validateAddComment(req *dto.AddMaintenanceRequestComment) error {
	if req.Content == "" {
		return app_errors.ErrCommentContentRequired
	}

	if len(req.Content) > 1000 {
		return app_errors.ErrCommentContentTooLong
	}

	return nil
}
