package app_errors

import "errors"

var (
	ErrRegisterNameRequired = errors.New("register request: name is required")
	ErrRegisterNameTooShort = errors.New("register request: name min length 2")
	ErrRegisterNameTooLong  = errors.New("register request: name max length 100")

	ErrRegisterEmailRequired = errors.New("register request: email is required")
	ErrRegisterEmailInvalid  = errors.New("register request: invalid email")

	ErrRegisterPasswordRequired = errors.New("register request: password is required")
	ErrRegisterPasswordTooShort = errors.New("register request: password min length 8")
	ErrRegisterPasswordTooLong  = errors.New("register request: password max length 72")
)

var (
	ErrLoginEmailRequired = errors.New("login request: email is required")
	ErrLoginEmailInvalid  = errors.New("login request: invalid email")
	ErrLoginPassRequired  = errors.New("login request: password is required")
)

var ErrRefreshTokenRequired = errors.New("refresh request: refresh_token is required")

var (
	ErrNewsTitleRequired = errors.New("create news request: title is required")
	ErrNewsTitleTooShort = errors.New("create news request: title min length 3")
	ErrNewsTitleTooLong  = errors.New("create news request: title max length 255")

	ErrNewsContentRequired = errors.New("create news request: content is required")
	ErrNewsContentTooShort = errors.New("create news request: content min length 10")
)

var (
	ErrRequestTitleRequired       = errors.New("create request: title is required")
	ErrRequestTitleTooShort       = errors.New("create request: title min length 3")
	ErrRequestTitleTooLong        = errors.New("create request: title max length 255")
	ErrRequestDescriptionRequired = errors.New("create request: description is required")
	ErrRequestDescriptionTooShort = errors.New("create request: description min length 10")
	ErrRequestTypeRequired        = errors.New("create request: type is required")
	ErrRequestTypeInvalid         = errors.New("create request: type must be one of: plumber, electrician")
)

var (
	ErrStatusRequired = errors.New("update status request: status is required")
	ErrStatusInvalid  = errors.New("update status request: status must be one of: open, in_progress, done, cancelled")
)

var (
	ErrCommentContentRequired = errors.New("add comment request: content is required")
	ErrCommentContentTooLong  = errors.New("add comment request: content max length 1000")
)

var (
	ErrDeviceTokenRequired = errors.New("register device request: device_token is required")
	ErrPlatformRequired    = errors.New("register device request: platform is required")
	ErrPlatformInvalid     = errors.New("register device request: platform must be one of: ios, android, web")
)
