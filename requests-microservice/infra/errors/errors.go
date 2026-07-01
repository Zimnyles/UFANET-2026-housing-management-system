package errors

import "errors"

var (
	ErrRequestNotFound = errors.New("request not found")
	ErrInternal        = errors.New("internal server error")
)

var (
	ErrRequestIDRequired       = errors.New("request id is required")
	ErrUserIDRequired          = errors.New("user_id is required")
	ErrTitleRequired           = errors.New("title is required")
	ErrTitleTooShort           = errors.New("title min length 3")
	ErrTitleTooLong            = errors.New("title max length 255")
	ErrDescriptionRequired     = errors.New("description is required")
	ErrDescriptionTooShort     = errors.New("description min length 10")
	ErrTypeRequired            = errors.New("type is required")
	ErrTypeInvalid             = errors.New("type must be one of: plumber, electrician")
	ErrStatusRequired          = errors.New("status is required")
	ErrStatusInvalid           = errors.New("status must be one of: open, in_progress, done, cancelled")
	ErrCommentContentRequired  = errors.New("comment content is required")
	ErrCommentContentTooLong   = errors.New("comment content max length 1000")
	ErrCommentRequestIDMissing = errors.New("comment request_id is required")
)
