package news_handler

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/dto"
)

func validateCreateNewsRequest(req *dto.CreateNewsRequest) error {
	if req.Title == "" {
		return app_errors.ErrNewsTitleRequired
	}
	if len(req.Title) < 3 {
		return app_errors.ErrNewsTitleTooShort
	}
	if len(req.Title) > 255 {
		return app_errors.ErrNewsTitleTooLong
	}
	if req.Content == "" {
		return app_errors.ErrNewsContentRequired
	}
	if len(req.Content) < 10 {
		return app_errors.ErrNewsContentTooShort
	}
	return nil
}
