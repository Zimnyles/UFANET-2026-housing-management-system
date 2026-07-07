package errors

import "errors"

var (
	ErrInternal              = errors.New("internal server error")
	ErrRegisterDevice        = errors.New("register device failed")
	ErrUnregisterDevice      = errors.New("unregister device failed")
	ErrListNotifications     = errors.New("list notifications failed")
	ErrDeliveryChannelClosed = errors.New("notification delivery channel closed")
)

var (
	ErrUserIDRequired      = errors.New("user_id is required")
	ErrDeviceTokenRequired = errors.New("device_token is required")
	ErrPlatformRequired    = errors.New("platform is required")
)
