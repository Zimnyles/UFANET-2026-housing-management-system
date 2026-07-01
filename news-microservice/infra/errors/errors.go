package errors

import "errors"

var (
	ErrNewsNotFound  = errors.New("news not found")
	ErrInternal      = errors.New("internal server error")
	ErrMQPublish     = errors.New("failed to publish message")
)

var (
	ErrTitleRequired   = errors.New("title is required")
	ErrTitleTooShort   = errors.New("title min length 3")
	ErrTitleTooLong    = errors.New("title max length 255")
	ErrContentRequired = errors.New("content is required")
	ErrContentTooShort = errors.New("content min length 10")
	ErrHouseIDRequired = errors.New("house_id is required")
	ErrNewsIDRequired  = errors.New("news id is required")
)
