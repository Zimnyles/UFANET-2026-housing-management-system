package server

import (
	newspb "github.com/zimnyles/UFANET-2026-housing-management-system/contracts/news/langs/go"

	infra_errors "news-service/infra/errors"
)

func ValidateCreateNewsRequest(req *newspb.CreateNewsRequest) error {
	if req.GetTitle() == "" {
		return infra_errors.ErrTitleRequired
	}

	if len(req.GetTitle()) < 3 {
		return infra_errors.ErrTitleTooShort
	}

	if len(req.GetTitle()) > 255 {
		return infra_errors.ErrTitleTooLong
	}

	if req.GetContent() == "" {
		return infra_errors.ErrContentRequired
	}

	if len(req.GetContent()) < 10 {
		return infra_errors.ErrContentTooShort
	}

	if req.GetHouseId() == "" {
		return infra_errors.ErrHouseIDRequired
	}

	return nil
}

func ValidateGetNewsItemRequest(req *newspb.GetNewsItemRequest) error {
	if req.GetId() == "" {
		return infra_errors.ErrNewsIDRequired
	}

	return nil
}
