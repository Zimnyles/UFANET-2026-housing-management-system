package server

import (
	requestspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/requests/langs/go"

	infra_errors "requests-service/infra/errors"
	"requests-service/infra/models/domain"
)

func ValidateCreateRequest(req *requestspb.CreateRequestRequest) error {
	if req.GetTitle() == "" {
		return infra_errors.ErrTitleRequired
	}

	if len(req.GetTitle()) < 3 {
		return infra_errors.ErrTitleTooShort
	}

	if len(req.GetTitle()) > 255 {
		return infra_errors.ErrTitleTooLong
	}

	if req.GetDescription() == "" {
		return infra_errors.ErrDescriptionRequired
	}

	if len(req.GetDescription()) < 10 {
		return infra_errors.ErrDescriptionTooShort
	}

	if req.GetType() == "" {
		return infra_errors.ErrTypeRequired
	}

	if !validType(req.GetType()) {
		return infra_errors.ErrTypeInvalid
	}

	if req.GetUserId() == "" {
		return infra_errors.ErrUserIDRequired
	}

	return nil
}

func ValidateGetRequest(req *requestspb.GetRequestRequest) error {
	if req.GetId() == "" {
		return infra_errors.ErrRequestIDRequired
	}

	return nil
}

func ValidateUpdateStatus(req *requestspb.UpdateStatusRequest) error {
	if req.GetId() == "" {
		return infra_errors.ErrRequestIDRequired
	}

	if req.GetStatus() == "" {
		return infra_errors.ErrStatusRequired
	}

	if !validStatus(req.GetStatus()) {
		return infra_errors.ErrStatusInvalid
	}

	if req.GetUserId() == "" {
		return infra_errors.ErrUserIDRequired
	}

	return nil
}

func ValidateAddComment(req *requestspb.AddCommentRequest) error {
	if req.GetRequestId() == "" {
		return infra_errors.ErrCommentRequestIDMissing
	}

	if req.GetUserId() == "" {
		return infra_errors.ErrUserIDRequired
	}

	if req.GetContent() == "" {
		return infra_errors.ErrCommentContentRequired
	}

	if len(req.GetContent()) > 1000 {
		return infra_errors.ErrCommentContentTooLong
	}

	return nil
}

func ValidateGetComments(req *requestspb.GetCommentsRequest) error {
	if req.GetRequestId() == "" {
		return infra_errors.ErrCommentRequestIDMissing
	}

	return nil
}

func validType(v string) bool {
	return v == domain.RequestTypePlumber || v == domain.RequestTypeElectrician
}

func validStatus(v string) bool {
	return v == domain.StatusOpen ||
		v == domain.StatusInProgress ||
		v == domain.StatusDone ||
		v == domain.StatusCancelled
}
