package notifications_handler

import (
	app_errors "api-gateway/internal/errors"
	"api-gateway/internal/models/dto"
)

func validateRegisterDeviceRequest(req *dto.RegisterDeviceRequest) error {
	if req.DeviceToken == "" {
		return app_errors.ErrDeviceTokenRequired
	}

	if req.Platform == "" {
		return app_errors.ErrPlatformRequired
	}

	if !req.Platform.Valid() {
		return app_errors.ErrPlatformInvalid
	}

	return nil
}
